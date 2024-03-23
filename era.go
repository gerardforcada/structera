package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type VersionedEraTemplateData struct {
	ExistingImports []string
	StructName      StructName
	Fields          []HubFieldInfo
	VersionNumber   int
}

func (g *Generator) EraFile(existingImports []string, version int, fields []HubFieldInfo) error {
	versionedDir := filepath.Join(g.OutputDir, g.Package, g.StructName.Snake)
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(versionedDir, fmt.Sprintf("v%d.go", version))); err == nil {
		if !g.Replace {
			fmt.Printf("Skipping existing versioned %s struct file: v%d.go\n", g.StructName.Original, version)
			return nil
		}
		fmt.Printf("Replacing existing versioned %s struct file: v%d.go\n", g.StructName.Original, version)
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "era.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, fmt.Sprintf("v%d.go", version)),
		Data: VersionedEraTemplateData{
			ExistingImports: existingImports,
			StructName:      g.StructName,
			Fields:          fields,
			VersionNumber:   version,
		},
	})
}
