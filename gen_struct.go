package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type VersionedStructTemplateData struct {
	PackageName     string
	ModulePackage   string
	ImportPath      string
	ExistingImports []string
	StructName      StructName
	VersionNumbers  []int
	VersionFields   map[int][]FieldInfo
}

func (g *Generator) StructFile(existingImports []string, importPath string) error {
	versionedDir := filepath.Join(g.OutputDir, "versioned")
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "templates/struct.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, fmt.Sprintf("%s.go", g.StructName.Lower)),
		Data: VersionedStructTemplateData{
			PackageName:     "versioned",
			ModulePackage:   string(ModulePackage),
			ImportPath:      importPath,
			ExistingImports: existingImports,
			StructName:      g.StructName,
			VersionNumbers:  g.Version.SortedVersions,
			VersionFields:   g.VersionedFields,
		},
	})
}
