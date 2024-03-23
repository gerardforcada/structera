package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_StructFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// Clean up after the test
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	type fields struct {
		StructName      StructName
		OutputDir       string
		Format          *Format
		VersionedFields map[int][]HubFieldInfo
		Package         string
		ProcessedFields []HubFieldInfo
	}
	tests := []struct {
		name            string
		fields          fields
		existingImports []string
		importPath      string
		wantErr         bool
		wantMatch       bool
	}{
		{
			name: "This config generates the example/versioned/testing.go file",
			fields: fields{
				StructName: StructName{
					Original: "Testing",
					Lower:    "testing",
					Snake:    "testing",
				},
				OutputDir: tempDir,
				Format: &Format{
					Versions: map[int][]string{
						1: {"InEveryVersion", "OnlyIn1", "FromStartTo3", "From1to4"},
						2: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						3: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						4: {"InEveryVersion", "From2ToEnd", "From1to4"},
					},
					SortedVersions: []int{1, 2, 3, 4},
				},
				VersionedFields: map[int][]HubFieldInfo{
					1: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "OnlyIn1",
							FormattedName: "OnlyIn1       ",
							Type:          "int",
							Tag:           "json:\"only_in_1\"",
						},
						{
							Name:          "FromStartTo3",
							FormattedName: "FromStartTo3  ",
							Type:          "[]byte",
							Tag:           "json:\"from_start_to_3\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
					2: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "From2ToEnd",
							FormattedName: "From2ToEnd    ",
							Type:          "uint8",
							Tag:           "json:\"from_2_to_end\"",
						},
						{
							Name:          "FromStartTo3",
							FormattedName: "FromStartTo3  ",
							Type:          "[]byte",
							Tag:           "json:\"from_start_to_3\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
					3: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "From2ToEnd",
							FormattedName: "From2ToEnd    ",
							Type:          "uint8",
							Tag:           "json:\"from_2_to_end\"",
						},
						{
							Name:          "FromStartTo3",
							FormattedName: "FromStartTo3  ",
							Type:          "[]byte",
							Tag:           "json:\"from_start_to_3\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
					4: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "From2ToEnd",
							FormattedName: "From2ToEnd    ",
							Type:          "uint8",
							Tag:           "json:\"from_2_to_end\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
				},
				Package: string(ModuleFolder),
				ProcessedFields: []HubFieldInfo{
					{FormattedName: "InEveryVersion", Type: "*string", Tag: "json:\"in_every_version\""},
					{FormattedName: "OnlyIn1", Type: "       *int", Tag: "json:\"only_in_1\""},
					{FormattedName: "From2ToEnd", Type: "    *uint8", Tag: "json:\"from_2_to_end\""},
					{FormattedName: "FromStartTo3", Type: "  *[]byte", Tag: "json:\"from_start_to_3\""},
					{FormattedName: "From1to4", Type: "      *float32", Tag: "json:\"from_1_to_4\""},
				},
			},
			existingImports: []string{},
			importPath:      "github.com/gerardforcada/structera/example",
			wantErr:         false,
			wantMatch:       true,
		},
		{
			name: "This config does not generate the example/versioned/testing.go file",
			fields: fields{
				StructName: StructName{
					Original: "Testing",
					Lower:    "testing",
					Snake:    "testing",
				},
				OutputDir: tempDir,
				Format: &Format{
					Versions: map[int][]string{
						1: {"InEveryVersion", "OnlyIn1", "FromStartTo3", "From1to4"},
						2: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						3: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						4: {"InEveryVersion", "From2ToEnd", "From1to4"},
					},
					SortedVersions: []int{1, 2, 3, 4},
				},
				VersionedFields: map[int][]HubFieldInfo{
					1: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "OnlyIn1",
							FormattedName: "OnlyIn1       ",
							Type:          "int",
							Tag:           "json:\"only_in_1\"",
						},
						{
							Name:          "FromStartTo3",
							FormattedName: "FromStartTo3  ",
							Type:          "[]byte",
							Tag:           "json:\"from_start_to_3\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
					2: {
						{
							Name:          "InEveryVersion",
							FormattedName: "InEveryVersion",
							Type:          "string",
							Tag:           "json:\"in_every_version\"",
						},
						{
							Name:          "From2ToEnd",
							FormattedName: "From2ToEnd    ",
							Type:          "uint8",
							Tag:           "json:\"from_2_to_end\"",
						},
						{
							Name:          "FromStartTo3",
							FormattedName: "FromStartTo3  ",
							Type:          "[]byte",
							Tag:           "json:\"from_start_to_3\"",
						},
						{
							Name:          "From1to4",
							FormattedName: "From1to4      ",
							Type:          "*float32",
							Tag:           "json:\"from_1_to_4\"",
						},
					},
				},
				Package: string(ModuleFolder),
				ProcessedFields: []HubFieldInfo{
					{FormattedName: "InEveryVersion", Type: "*string"},
					{FormattedName: "OnlyIn1", Type: "*int"},
					{FormattedName: "FromStartTo3", Type: "[]byte"},
					{FormattedName: "From1to4", Type: "*float32"},
					{FormattedName: "From2ToEnd", Type: "*uint8"},
				},
			},
			existingImports: []string{"test"},
			importPath:      "github.com/gerardforcada/structera/example",
			wantErr:         false,
			wantMatch:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				StructName:      tt.fields.StructName,
				OutputDir:       tt.fields.OutputDir,
				Format:          tt.fields.Format,
				VersionedFields: tt.fields.VersionedFields,
				Package:         tt.fields.Package,
				ProcessedFields: tt.fields.ProcessedFields,
			}

			if err := g.HubFile(tt.existingImports, tt.importPath); (err != nil) != tt.wantErr {
				t.Errorf("HubFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Path to the generated file
			generatedFilePath := filepath.Join(tempDir, tt.fields.Package, fmt.Sprintf("%s.go", tt.fields.StructName.Lower))

			// Read the contents of the generated file
			generatedFileContent, err := os.ReadFile(generatedFilePath)
			assert.NoError(t, err)

			// Path to the reference file
			referenceFilePath := filepath.Join("example", tt.fields.Package, fmt.Sprintf("%s.go", tt.fields.StructName.Lower))

			// Read the contents of the reference file
			referenceFileContent, err := os.ReadFile(referenceFilePath)
			assert.NoError(t, err)

			// Compare the contents
			if tt.wantMatch {
				assert.Equal(t, string(referenceFileContent), string(generatedFileContent), "The generated file content does not match the reference file content")
			} else {
				assert.NotEqual(t, string(referenceFileContent), string(generatedFileContent), "The generated file content does not match the reference file content")
			}
		})
	}
}
