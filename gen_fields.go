package main

import (
	"os"
	"path/filepath"
)

type FieldInfo struct {
	Name          string
	FormattedName string
	Type          string
	Tag           string // This can be empty if there's no tag
}

type FieldsTemplateData struct {
	PackageName string
	Fields      []FieldInfo
}

func (g *Generator) FieldsFile() error {
	fieldsDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(fieldsDir, os.ModePerm); err != nil {
		return err
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "fields.go.tmpl",
		OutputFilePath:   filepath.Join(fieldsDir, "fields.go"),
		Data: FieldsTemplateData{
			PackageName: g.StructName.Lower,
			Fields:      g.ProcessedFields,
		},
	})
}
