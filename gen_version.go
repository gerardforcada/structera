package main

import (
	"os"
	"path/filepath"
)

type VersionTemplateData struct {
	PackageName   string
	ModulePackage string
	Versions      []int
}

func (g *Generator) VersionFile(versions []int) error {
	versionDir := filepath.Join(g.OutputDir, "versioned", g.StructName.Lower)
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "version.go.tmpl",
		OutputFilePath:   filepath.Join(versionDir, "version.go"),
		Data: VersionTemplateData{
			PackageName:   g.StructName.Lower,
			ModulePackage: string(ModulePackage),
			Versions:      versions,
		},
	})
}
