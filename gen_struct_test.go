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
		Version         *Version
		VersionedFields map[int][]FieldInfo
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
				},
				OutputDir: tempDir,
				Version: &Version{
					Versions: map[int][]string{
						1: {"InEveryVersion", "OnlyIn1", "FromStartTo3", "From1to4"},
						2: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						3: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						4: {"InEveryVersion", "From2ToEnd", "From1to4"},
					},
					SortedVersions: []int{1, 2, 3, 4},
				},
				VersionedFields: map[int][]FieldInfo{
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
			},
			existingImports: []string{},
			importPath:      "github.com/gerardforcada/structera/example/versioned/testing",
			wantErr:         false,
			wantMatch:       true,
		},
		{
			name: "This config does not generate the example/versioned/testing.go file",
			fields: fields{
				StructName: StructName{
					Original: "Testing",
					Lower:    "testing",
				},
				OutputDir: tempDir,
				Version: &Version{
					Versions: map[int][]string{
						1: {"InEveryVersion", "OnlyIn1", "FromStartTo3", "From1to4"},
						2: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						3: {"InEveryVersion", "From2ToEnd", "FromStartTo3", "From1to4"},
						4: {"InEveryVersion", "From2ToEnd", "From1to4"},
					},
					SortedVersions: []int{1, 2, 3, 4},
				},
				VersionedFields: map[int][]FieldInfo{
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
			},
			existingImports: []string{"test"},
			importPath:      "github.com/gerardforcada/structera/example/versioned/testing",
			wantErr:         false,
			wantMatch:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				StructName:      tt.fields.StructName,
				OutputDir:       tt.fields.OutputDir,
				Version:         tt.fields.Version,
				VersionedFields: tt.fields.VersionedFields,
			}
			if err := g.StructFile(tt.existingImports, tt.importPath); (err != nil) != tt.wantErr {
				t.Errorf("StructFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Path to the generated file
			generatedFilePath := filepath.Join(tempDir, "versioned", fmt.Sprintf("%s.go", tt.fields.StructName.Lower))

			// Read the contents of the generated file
			generatedFileContent, err := os.ReadFile(generatedFilePath)
			assert.NoError(t, err)

			// Path to the reference file
			referenceFilePath := filepath.Join("example", "versioned", fmt.Sprintf("%s.go", tt.fields.StructName.Lower))

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
