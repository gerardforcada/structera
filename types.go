package main

import (
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
)

type VersionedTypesTemplateData struct {
	PackageName   string
	ModulePackage string
	FileNameMap   map[string]string
}

func (g *Generator) TypesFile(importPath string) error {
	versionedDir := filepath.Join(g.OutputDir, g.Package)
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}

	files, err := os.ReadDir(filepath.Join(g.OutputDir, g.Package))
	if err != nil {
		return err
	}

	fileNameMap := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if filepath.Ext(fileName) == ".go" && fileName != "types.go" {
				snakeCaseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
				camelCaseName := strcase.UpperCamelCase(snakeCaseName)
				fileNameMap[camelCaseName] = snakeCaseName
			}
		}
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "types.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, "types.go"),
		Data: VersionedTypesTemplateData{
			PackageName:   g.Package,
			ModulePackage: string(ModulePackage),
			FileNameMap:   fileNameMap,
		},
	})
}
