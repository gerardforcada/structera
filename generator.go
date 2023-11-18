package main

import (
	"fmt"
	"github.com/gerardforcada/structera/helpers"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type VersionTemplateData struct {
	PackageName   string
	ModulePackage string
	Versions      []int
}

type FieldInfo struct {
	Name          string
	FormattedName string
	Type          string
	Tag           string // This can be empty if there's no tag
}

func (fi FieldInfo) FormatField() string {
	if fi.Tag != "" {
		return fmt.Sprintf("%s %s `%s`", fi.FormattedName, fi.Type, fi.Tag)
	}
	return fmt.Sprintf("%s %s", fi.FormattedName, fi.Type)
}

type FieldsTemplateData struct {
	PackageName string
	Fields      []FieldInfo
}

type TemplateData any

type GenerateFileFromTemplateInput struct {
	TemplateFilePath string
	OutputFilePath   string
	Data             TemplateData
}

type VersionedStructTemplateData struct {
	PackageName     string
	ModulePackage   string
	ImportPath      string
	ExistingImports []string
	StructName      StructName // Assuming StructName is a struct with appropriate fields
	VersionNumbers  []int
	VersionFields   map[int][]FieldInfo // Map of version number to slice of FieldInfo
}

func GenerateFileFromTemplate(input GenerateFileFromTemplateInput) error {
	// Ensure the directory for the output file exists
	outputDir := filepath.Dir(input.OutputFilePath)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	// Parse the template file
	tmpl, err := template.New(filepath.Base(input.TemplateFilePath)).Funcs(template.FuncMap{"sub": helpers.Sub}).ParseFiles(input.TemplateFilePath)
	if err != nil {
		return err
	}

	// Create the output file
	file, err := os.Create(input.OutputFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	// Execute the template with the data
	return tmpl.Execute(file, input.Data)
}

func (g Generator) GenerateVersionFile(versions []int) error {
	versionDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	return GenerateFileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "templates/version.go.tmpl",
		OutputFilePath:   filepath.Join(versionDir, "version.go"),
		Data: VersionTemplateData{
			PackageName:   g.StructName.Lower,
			ModulePackage: string(ModulePackage),
			Versions:      versions,
		},
	})
}

func (g Generator) GenerateFieldsFile(structType *ast.StructType) error {
	fieldsDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(fieldsDir, os.ModePerm); err != nil {
		return err
	}

	fields, maxNameLength, err := g.processFieldInfo(structType)
	if err != nil {
		return err
	}

	// Second pass to add padding for alignment
	for i := range fields {
		padding := maxNameLength - len(fields[i].Name)
		fields[i].FormattedName = fields[i].Name + strings.Repeat(" ", padding)
	}

	return GenerateFileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "templates/fields.go.tmpl",
		OutputFilePath:   filepath.Join(fieldsDir, "fields.go"),
		Data: FieldsTemplateData{
			PackageName: g.StructName.Lower,
			Fields:      fields,
		},
	})
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

			err = g.GenerateStructFile(structType, imports, importPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("struct '%s' not found in file '%s'", g.StructName.Original, g.Filename)
}

func (g Generator) GenerateStructFile(structType *ast.StructType, existingImports []string, importPath string) error {
	versionedDir := filepath.Join(g.OutputDir, "versioned")
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}

	versionFields, maxNameLength, err := g.prepareVersionFields(structType)
	if err != nil {
		return err
	}

	// Apply padding for alignment to version-specific fields
	for _, fields := range versionFields {
		for i := range fields {
			padding := maxNameLength - len(fields[i].Name)
			fields[i].FormattedName = fields[i].Name + strings.Repeat(" ", padding)
		}
	}

	return GenerateFileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "templates/struct.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, fmt.Sprintf("%s.go", g.StructName.Lower)),
		Data: VersionedStructTemplateData{
			PackageName:     "versioned",
			ModulePackage:   string(ModulePackage),
			ImportPath:      importPath,
			ExistingImports: existingImports,
			StructName:      g.StructName,
			VersionNumbers:  g.Version.SortedVersions,
			VersionFields:   versionFields,
		},
	})
}

func (g Generator) prepareVersionFields(structType *ast.StructType) (map[int][]FieldInfo, int, error) {
	fields, maxNameLength, err := g.processFieldInfo(structType)
	if err != nil {
		return nil, 0, err
	}

	versionFields := make(map[int][]FieldInfo)
	for version, versionedFieldStrs := range g.Version.Versions {
		var versionFieldInfos []FieldInfo

		for _, versionedFieldStr := range versionedFieldStrs {
			parts := strings.SplitN(versionedFieldStr, " ", 2)
			if len(parts) < 2 {
				continue
			}
			fieldName := parts[0]

			for _, field := range fields {
				if field.Name == fieldName {
					versionFieldInfos = append(versionFieldInfos, field)
					break
				}
			}
		}

		versionFields[version] = versionFieldInfos
	}

	return versionFields, maxNameLength, nil
}

func (g Generator) processFieldInfo(structType *ast.StructType) ([]FieldInfo, int, error) {
	var fields []FieldInfo
	maxNameLength := 0

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		fieldType := formatFieldType(field.Type, true) // Assuming this function formats the type as string

		fieldInfo := FieldInfo{
			Name: fieldName,
			Type: fieldType,
		}

		if field.Tag != nil {
			tagValue := field.Tag.Value
			tag := tagValue[1 : len(tagValue)-1] // Extract tag string without quotes
			filteredTag := g.Version.ExcludeVersionTag(tag)
			if filteredTag != "" {
				fieldInfo.Tag = filteredTag
			}
		}

		fields = append(fields, fieldInfo)

		if len(fieldName) > maxNameLength {
			maxNameLength = len(fieldName)
		}
	}

	return fields, maxNameLength, nil
}
