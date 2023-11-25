package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type FieldInfo struct {
	Name          string
	FormattedName string
	Type          string
	Tag           string
}

type VersionedStructTemplateData struct {
	PackageName     string
	ModulePackage   string
	ImportPath      string
	ExistingImports []string
	StructName      StructName
	Fields          []FieldInfo
	VersionedFields map[int][]FieldInfo
	Versions        []int
	CustomType      bool
}

func (g *Generator) StructFile(existingImports []string, importPath string) error {
	versionedDir := filepath.Join(g.OutputDir, g.Package)
	if err := os.MkdirAll(versionedDir, os.ModePerm); err != nil {
		return err
	}

	return g.FileFromTemplate(GenerateFileFromTemplateInput{
		TemplateFilePath: "struct.go.tmpl",
		OutputFilePath:   filepath.Join(versionedDir, fmt.Sprintf("%s.go", g.StructName.Snake)),
		Data: VersionedStructTemplateData{
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
