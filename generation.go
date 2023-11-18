package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func formatStructFields(fields []string, tags map[string]string) string {
	maxLen := 0
	for _, f := range fields {
		parts := strings.Split(f, " ")
		if len(parts[0]) > maxLen {
			maxLen = len(parts[0])
		}
	}

	var buf strings.Builder
	for _, f := range fields {
		parts := strings.Split(f, " ")
		fieldName := parts[0]
		fieldType := parts[1]
		buf.WriteString(fmt.Sprintf("\t%s%s %s", fieldName, strings.Repeat(" ", maxLen-len(fieldName)), fieldType))
		if tag, ok := tags[fieldName]; ok {
			buf.WriteString(fmt.Sprintf(" `%s`", tag))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func excludeVersionTag(tag string) string {
	var result []string
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		if !strings.HasPrefix(t, "version:") {
			result = append(result, t)
		}
	}
	return strings.Join(result, " ")
}

func generateVersionFile(structName string, versions []int, outputDir string) error {
	versionDir := filepath.Join(outputDir, "versioned", strings.ToLower(structName))
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(structName)))
	buf.WriteString(fmt.Sprintf("import \"%s/version\"\n\n", LibraryPackage))
	buf.WriteString("type Version version.Version\n\nconst (\n")

	for _, v := range versions {
		buf.WriteString(fmt.Sprintf("\tVersion%d Version = %d\n", v, v))
	}

	buf.WriteString(")\n\n")
	buf.WriteString("type Versioned[T any] version.Versioned[T]\n\n")

	for _, v := range versions {
		buf.WriteString(fmt.Sprintf("type V%d[T Versioned[Version]] *T\n", v))
	}

	versionFileName := filepath.Join(versionDir, "version.go")
	return os.WriteFile(versionFileName, []byte(buf.String()), 0644)
}

func generateFieldsFile(structName string, structType *ast.StructType, outputDir string) error {
	fieldsDir := filepath.Join(outputDir, "versioned", strings.ToLower(structName))
	if err := os.MkdirAll(fieldsDir, os.ModePerm); err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(structName)))
	buf.WriteString("type PointerFields struct {\n")

	var fields []string
	tags := make(map[string]string)
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fieldName := name.Name
			fieldType := formatFieldType(field.Type, true)
			fields = append(fields, fmt.Sprintf("%s %s", fieldName, fieldType))
			if field.Tag != nil {
				tagValue := field.Tag.Value
				tag := tagValue[1 : len(tagValue)-1] // Extract tag string without quotes
				filteredTag := excludeVersionTag(tag)
				if filteredTag != "" {
					tags[fieldName] = filteredTag
				}
			}
		}
	}

	buf.WriteString(formatStructFields(fields, tags))
	buf.WriteString("}\n")

	fieldsFileName := filepath.Join(fieldsDir, "fields.go")
	return os.WriteFile(fieldsFileName, []byte(buf.String()), 0644)
}

func formatFieldType(expr ast.Expr, pointer bool) string {
	result := ""
	if pointer {
		result += "*"
	}

	switch t := expr.(type) {
	case *ast.Ident:
		result += t.Name // Simple types like int, string, etc.
	case *ast.ArrayType:
		// For arrays and slices
		elementType := formatFieldType(t.Elt, false) // No extra pointer for element type
		result += fmt.Sprintf("[]%s", elementType)
	case *ast.MapType:
		// For maps
		keyType := formatFieldType(t.Key, false)     // No extra pointer for key type
		valueType := formatFieldType(t.Value, false) // No extra pointer for value type
		result += fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.StarExpr:
		// For pointers, simply add an extra star
		pointedType := formatFieldType(t.X, false)
		result += fmt.Sprintf("%s", pointedType)
	case *ast.StructType:
		// Handle nested structs
		result += "NestedStruct" // Adjust as needed for nested structs
	case *ast.SelectorExpr:
		// For qualified identifiers (e.g., time.Time)
		result += fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	default:
		// Fallback for other types
		result += "interface{}" // Generic pointer type as fallback
	}

	return result
}

func generateMethods(structName string, lowerCaseStructName string, versionNumbers []int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func (d %s) GetVersionStructs() []version.Versioned[%s.Version] {\n", structName, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn []version.Versioned[%s.Version]{\n", lowerCaseStructName))
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\t%sV%d{},\n", structName, v))
	}
	buf.WriteString("\t}\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) GetBaseStruct() interface{} {\n", structName))
	buf.WriteString("\treturn d.PointerFields\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) DetectVersion() %s.Version {\n", structName, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn version.DetectBestMatch[%s.Version, %s](d)\n}\n\n", lowerCaseStructName, structName))

	buf.WriteString(fmt.Sprintf("func (d %s) GetVersions() []%s.Version {\n", structName, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn []%s.Version{\n", lowerCaseStructName))
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\t%s.Version%d,\n", lowerCaseStructName, v))
	}
	buf.WriteString("\t}\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) GetMinVersion() %s.Version {\n", structName, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version1\n}\n\n", lowerCaseStructName))

	maxVersion := versionNumbers[len(versionNumbers)-1]
	buf.WriteString(fmt.Sprintf("func (d %s) GetMaxVersion() %s.Version {\n", structName, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version%d\n}\n\n", lowerCaseStructName, maxVersion))

	return buf.String()
}

func generateVersionMethods(structName, lowerCaseStructName string, version int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func (d %sV%d) GetVersion() %s.Version {\n", structName, version, lowerCaseStructName))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version%d\n}\n\n", lowerCaseStructName, version))

	return buf.String()
}

func generateVersionedStructs(fileName, structName, outputDir string) error {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var imports []string
	for _, i := range node.Imports {
		imports = append(imports, strings.Trim(i.Path.Value, "\""))
	}

	lowerCaseStructName := strings.ToLower(structName)

	goModPath, err := findGoModPath(filepath.Dir(fileName))
	if err != nil {
		return fmt.Errorf("error finding go.mod: %v", err)
	}

	baseImportPath, err := getBaseImportPath(goModPath)
	if err != nil {
		return fmt.Errorf("error getting base import path: %v", err)
	}

	relativePath, err := filepath.Rel(outputDir, filepath.Dir(fileName))
	if err != nil {
		return err
	}
	importPath := path.Join(baseImportPath, relativePath, "versioned", lowerCaseStructName)

	for _, f := range node.Decls {
		genDecl, ok := f.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			versions := identifyVersions(structType)
			if len(versions) == 0 {
				return fmt.Errorf("no version tags found in struct")
			}

			// Generate fields.go
			if err := generateFieldsFile(structName, structType, outputDir); err != nil {
				return err
			}

			// Generate version.go
			var versionNumbers []int
			for v := range versions {
				versionNumbers = append(versionNumbers, v)
			}
			sort.Ints(versionNumbers)
			if err := generateVersionFile(structName, versionNumbers, outputDir); err != nil {
				return err
			}

			versionedStructsCode, err := processStruct(structName, lowerCaseStructName, structType, imports, importPath)
			if err != nil {
				return err
			}

			// Create a versioned folder
			versionedDir := filepath.Join(outputDir, "versioned")
			if _, err := os.Stat(versionedDir); os.IsNotExist(err) {
				err := os.MkdirAll(versionedDir, os.ModePerm)
				if err != nil {
					return err
				}
			}

			outputFileName := filepath.Join(versionedDir, strings.ToLower(structName)+".go")
			return os.WriteFile(outputFileName, []byte(versionedStructsCode), 0644)
		}
	}

	return fmt.Errorf("struct '%s' not found in file '%s'", structName, fileName)
}
