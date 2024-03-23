package main

import (
	"fmt"
	"github.com/gerardforcada/structera/helpers"
	"github.com/gerardforcada/structera/templates"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type StructName struct {
	Original string
	Lower    string
	Snake    string
}

type Generator struct {
	Format          *Format
	Resolver        *Resolver
	Filename        string
	StructName      StructName
	OutputDir       string
	ProcessedFields []HubFieldInfo
	VersionedFields map[int][]HubFieldInfo
	Package         string
	Replace         bool
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

	tmpl, err := template.New(filepath.Base(input.TemplateFilePath)).Funcs(template.FuncMap{"sub": helpers.Sub}).ParseFS(templates.FS, input.TemplateFilePath)
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

	err = g.Resolver.GetBaseImportPath()
	if err != nil {
		return fmt.Errorf("error getting base import path: %v", err)
	}

	// Calculate the relative path from the go.mod directory to the generated file directory
	goModDir := filepath.Dir(g.Resolver.GoModPath)
	relativePath, err := filepath.Rel(goModDir, g.OutputDir)
	if err != nil {
		return err
	}

	importPath := path.Join(g.Resolver.ImportPath, relativePath)

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

			g.Format.IdentifyVersions(structType)
			if len(g.Format.Versions) == 0 {
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
			g.PrepareVersionedFields()

			// Generate versioned struct files
			err = g.HubFile(imports, importPath)
			if err != nil {
				return err
			}

			for version, fields := range g.VersionedFields {
				err = g.EraFile(imports, version, fields)
				if err != nil {
					return err
				}
			}
			// Generate types.go file
			err = g.TypesFile(importPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("struct '%s' not found in file '%s'", g.StructName.Original, g.Filename)
}

func (g *Generator) PrepareVersionedFields() {
	versionedFields := make(map[int][]HubFieldInfo)
	for version, versionedFieldStrings := range g.Format.Versions {
		var versionFieldInfos []HubFieldInfo

		for _, versionedFieldStr := range versionedFieldStrings {
			parts := strings.SplitN(versionedFieldStr, " ", 2)
			if len(parts) < 2 {
				continue
			}
			fieldName := parts[0]

			for _, field := range g.ProcessedFields {
				if field.Name == fieldName {
					field.Type = field.Type[1:] // remove first char of field type (asterisk)
					versionFieldInfos = append(versionFieldInfos, field)
					break
				}
			}
		}

		versionedFields[version] = versionFieldInfos
	}

	g.VersionedFields = versionedFields
}

func (g *Generator) ProcessFieldInfo(structType *ast.StructType) ([]HubFieldInfo, int, error) {
	var fields []HubFieldInfo
	maxNameLength := 0

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		fieldType := g.Format.FieldType(field.Type, true)

		fieldInfo := HubFieldInfo{
			Name: fieldName,
			Type: fieldType,
		}

		if field.Tag != nil {
			tagValue := field.Tag.Value
			tag := tagValue[1 : len(tagValue)-1] // Extract tag string without quotes
			filteredTag := g.Format.ExcludeVersionTag(tag)
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
