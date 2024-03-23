package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_EraFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// Clean up after the test
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	type fields struct {
		ExistingImports []string
		StructName      StructName
		Fields          []HubFieldInfo
		VersionNumber   int
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
				Fields: []HubFieldInfo{
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
			},
			existingImports: []string{},
			importPath:      "github.com/gerardforcada/structera/example",
			wantErr:         false,
			wantMatch:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				StructName:      tt.fields.StructName,
				OutputDir:       tempDir,
				Format:          &Format{},
				VersionedFields: map[int][]HubFieldInfo{},
			}

			if err := g.EraFile(tt.existingImports, 1, tt.fields.Fields); (err != nil) != tt.wantErr {
				t.Errorf("HubFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Path to the generated file
			generatedFilePath := filepath.Join(tempDir, tt.fields.StructName.Lower, fmt.Sprintf("v%d.go", 1))

			// Read the contents of the generated file
			generatedFileContent, err := os.ReadFile(generatedFilePath)
			assert.NoError(t, err)

			// Path to the reference file
			referenceFilePath := filepath.Join("example", string(ModuleFolder), tt.fields.StructName.Lower, fmt.Sprintf("v%d.go", 1))

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
