package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type HubFieldInfo struct {
	Name          string
	FormattedName string
	Type          string
	Tag           string
}

type VersionedHubTemplateData struct {
	PackageName     string
	ModulePackage   string
	ImportPath      string
	ExistingImports []string
	StructName      StructName
	Fields          []HubFieldInfo
	VersionedFields map[int][]HubFieldInfo
	Versions        []int
	CustomType      bool
}

func (g *Generator) HubFile(existingImports []string, importPath string) error {
	versionedDir := filepath.Join(g.OutputDir, g.Package)
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "hub.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, fmt.Sprintf("%s.go", g.StructName.Snake)),
		Data: VersionedHubTemplateData{
			PackageName:     g.Package,
			ModulePackage:   string(ModulePackage),
			ImportPath:      importPath,
			ExistingImports: existingImports,
			StructName:      g.StructName,
			VersionedFields: g.VersionedFields,
			Fields:          g.ProcessedFields,
			Versions:        g.Format.SortedVersions,
			CustomType:      g.Format.CustomType,
		},
	})
}
