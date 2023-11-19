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

type StructName struct {
	Original string
	Lower    string
}

type Generator struct {
	Version         *Version
	Resolver        *Resolver
	Filename        string
	StructName      StructName
	OutputDir       string
	ProcessedFields []FieldInfo
	VersionedFields map[int][]FieldInfo
}

type GenerateFileFromTemplateInput struct {
	TemplateFilePath string
	OutputFilePath   string
	Data             any
}

func (g *Generator) FileFromTemplate(input GenerateFileFromTemplateInput) error {
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

func (g *Generator) VersionedStructs() error {
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

			fields, maxNameLength, err := g.ProcessFieldInfo(structType)
			if err != nil {
				return err
			}

			// Second pass to add padding for alignment
			for i := range fields {
				padding := maxNameLength - len(fields[i].Name)
				fields[i].FormattedName = fields[i].Name + strings.Repeat(" ", padding)
			}

			g.ProcessedFields = fields
			g.PrepareVersionFields()

			// Generate fields.go
			if err := g.FieldsFile(); err != nil {
				return err
			}

			// Generate version.go
			var versionNumbers []int
			for v := range g.Version.Versions {
				versionNumbers = append(versionNumbers, v)
			}
			sort.Ints(versionNumbers)
			if err := g.VersionFile(versionNumbers); err != nil {
				return err
			}

			err = g.StructFile(imports, importPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("struct '%s' not found in file '%s'", g.StructName.Original, g.Filename)
}

func (g *Generator) PrepareVersionFields() {
	versionFields := make(map[int][]FieldInfo)
	for version, versionedFieldStrs := range g.Version.Versions {
		var versionFieldInfos []FieldInfo

		for _, versionedFieldStr := range versionedFieldStrs {
			parts := strings.SplitN(versionedFieldStr, " ", 2)
			if len(parts) < 2 {
				continue
			}
			fieldName := parts[0]

			for _, field := range g.ProcessedFields {
				if field.Name == fieldName {
					versionFieldInfos = append(versionFieldInfos, field)
					break
				}
			}
		}

		versionFields[version] = versionFieldInfos
	}

	g.VersionedFields = versionFields
}

func (g *Generator) ProcessFieldInfo(structType *ast.StructType) ([]FieldInfo, int, error) {
	var fields []FieldInfo
	maxNameLength := 0

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		fieldType := formatFieldType(field.Type, true)

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
