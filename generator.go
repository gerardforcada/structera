package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func (g Generator) GenerateVersionFile(versions []int) error {
	versionDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("import \"%s/version\"\n\n", ModulePackage))
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

func (g Generator) GenerateFieldsFile(structType *ast.StructType) error {
	fieldsDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(fieldsDir, os.ModePerm); err != nil {
		return err
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.StructName.Lower))
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
				filteredTag := g.Version.ExcludeVersionTag(tag)
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

func (g Generator) GenerateMethods(versionNumbers []int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func (d %s) GetVersionStructs() []version.Versioned[%s.Version] {\n", g.StructName.Original, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn []version.Versioned[%s.Version]{\n", g.StructName.Lower))
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\t%sV%d{},\n", g.StructName.Original, v))
	}
	buf.WriteString("\t}\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) GetBaseStruct() interface{} {\n", g.StructName.Original))
	buf.WriteString("\treturn d.PointerFields\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) DetectVersion() %s.Version {\n", g.StructName.Original, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn version.DetectBestMatch[%s.Version, %s](d)\n}\n\n", g.StructName.Lower, g.StructName.Original))

	buf.WriteString(fmt.Sprintf("func (d %s) GetVersions() []%s.Version {\n", g.StructName.Original, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn []%s.Version{\n", g.StructName.Lower))
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\t%s.Version%d,\n", g.StructName.Lower, v))
	}
	buf.WriteString("\t}\n}\n\n")

	buf.WriteString(fmt.Sprintf("func (d %s) GetMinVersion() %s.Version {\n", g.StructName.Original, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version1\n}\n\n", g.StructName.Lower))

	maxVersion := versionNumbers[len(versionNumbers)-1]
	buf.WriteString(fmt.Sprintf("func (d %s) GetMaxVersion() %s.Version {\n", g.StructName.Original, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version%d\n}\n\n", g.StructName.Lower, maxVersion))

	return buf.String()
}

func (g Generator) GenerateVersionMethods(version int) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func (d %sV%d) GetVersion() %s.Version {\n", g.StructName.Original, version, g.StructName.Lower))
	buf.WriteString(fmt.Sprintf("\treturn %s.Version%d\n}\n\n", g.StructName.Lower, version))

	return buf.String()
}

func (g Generator) GenerateVersionedStructs() error {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, g.Filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var imports []string
	for _, i := range node.Imports {
		imports = append(imports, strings.Trim(i.Path.Value, "\""))
	}

	err = g.Resolver.FindGoModPath(filepath.Dir(g.Filename))
	if err != nil {
		return fmt.Errorf("error finding go.mod: %v", err)
	}

	err = g.Resolver.GetBaseImportPath(g.Resolver.GoModPath)
	if err != nil {
		return fmt.Errorf("error getting base import path: %v", err)
	}

	relativePath, err := filepath.Rel(g.OutputDir, filepath.Dir(g.Filename))
	if err != nil {
		return err
	}
	importPath := path.Join(g.Resolver.ImportPath, relativePath, "versioned", g.StructName.Lower)

	for _, f := range node.Decls {
		genDecl, ok := f.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != g.StructName.Original {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			g.Version.IdentifyVersions(structType)
			if len(g.Version.Versions) == 0 {
				return fmt.Errorf("no version tags found in struct")
			}

			// Generate fields.go
			if err := g.GenerateFieldsFile(structType); err != nil {
				return err
			}

			// Generate version.go
			var versionNumbers []int
			for v := range g.Version.Versions {
				versionNumbers = append(versionNumbers, v)
			}
			sort.Ints(versionNumbers)
			if err := g.GenerateVersionFile(versionNumbers); err != nil {
				return err
			}

			versionedStructsCode, err := g.GenerateStructFile(structType, imports, importPath)
			if err != nil {
				return err
			}

			// Create a versioned folder
			versionedDir := filepath.Join(g.OutputDir, "versioned")
			if _, err := os.Stat(versionedDir); os.IsNotExist(err) {
				err := os.MkdirAll(versionedDir, os.ModePerm)
				if err != nil {
					return err
				}
			}

			outputFileName := filepath.Join(versionedDir, strings.ToLower(g.StructName.Original)+".go")
			return os.WriteFile(outputFileName, []byte(versionedStructsCode), 0644)
		}
	}

	return fmt.Errorf("struct '%s' not found in file '%s'", g.StructName.Original, g.Filename)
}

func (g Generator) GenerateStructFile(structType *ast.StructType, existingImports []string, importPath string) (string, error) {
	// Sort the version numbers
	var versionNumbers []int
	for v := range g.Version.Versions {
		versionNumbers = append(versionNumbers, v)
	}
	sort.Ints(versionNumbers)

	var buf strings.Builder
	buf.WriteString("package versioned\n\n")
	buf.WriteString("import (\n")
	buf.WriteString(fmt.Sprintf("\t\"%s/version\"\n", ModulePackage))
	buf.WriteString(fmt.Sprintf("\t\"%s\"\n", importPath)) // Add dynamic import path

	// Include existing imports from the original file
	for _, imp := range existingImports {
		buf.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}
	buf.WriteString(")\n\n")

	// Versions struct
	buf.WriteString(fmt.Sprintf("type %sVersions struct {\n", g.StructName.Original))
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\tV%d %s.V%d[%sV%d]\n", v, g.StructName.Lower, v, g.StructName.Original, v))
	}
	buf.WriteString("}\n\n")

	// struct
	buf.WriteString(fmt.Sprintf("type %s struct {\n\t%s.PointerFields\n\t%sVersions\n}\n\n", g.StructName.Original, g.StructName.Lower, g.StructName.Original))

	// Initialize function
	buf.WriteString(fmt.Sprintf("func (d *%s) Initialize() {\n", g.StructName.Original))
	buf.WriteString(fmt.Sprintf("\td.%sVersions = %sVersions{\n", g.StructName.Original, g.StructName.Original))

	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\tV%d: &%sV%d{},\n", v, g.StructName.Original, v))
	}

	buf.WriteString("\t}\n}\n\n")

	// Additional methods for the struct
	buf.WriteString(g.GenerateMethods(versionNumbers))

	// Version-specific struct types and methods
	for _, v := range versionNumbers {
		fields, ok := g.Version.Versions[v]
		if !ok {
			continue // Skip if version number is not found in the map
		}

		// Extract tags from original struct fields
		tags := make(map[string]string)
		for _, field := range structType.Fields.List {
			if field.Tag != nil {
				for _, name := range field.Names {
					tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json")
					if tag != "" {
						tags[name.Name] = fmt.Sprintf("json:\"%s\"", tag)
					}
				}
			}
		}

		buf.WriteString(fmt.Sprintf("type %sV%d struct {\n", g.StructName.Original, v))
		buf.WriteString(formatStructFields(fields, tags))
		buf.WriteString("}\n\n")
		buf.WriteString(g.GenerateVersionMethods(v))
	}

	return buf.String(), nil
}
