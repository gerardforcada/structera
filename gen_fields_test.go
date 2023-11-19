package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_FieldsFile(t *testing.T) {

	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// Clean up after the test
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	type fields struct {
		Version         *Version
		Resolver        *Resolver
		Filename        string
		StructName      StructName
		OutputDir       string
		ProcessedFields []FieldInfo
		VersionedFields map[int][]FieldInfo
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantMatch bool
	}{
		{
			name: "This config generates the example/versioned/testing/fields.go file",
			fields: fields{
				StructName: StructName{
					Lower: "testing",
				},
				OutputDir: tempDir,
				ProcessedFields: []FieldInfo{
					{
						Name:          "InEveryVersion",
						FormattedName: "InEveryVersion",
						Type:          "*string",
						Tag:           "json:\"in_every_version\"",
					},
					{
						Name:          "OnlyIn1",
						FormattedName: "OnlyIn1       ",
						Type:          "*int",
						Tag:           "json:\"only_in_1\"",
					},
					{
						Name:          "From2ToEnd",
						FormattedName: "From2ToEnd    ",
						Type:          "*uint8",
						Tag:           "json:\"from_2_to_end\"",
					},
					{
						Name:          "FromStartTo3",
						FormattedName: "FromStartTo3  ",
						Type:          "*[]byte",
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
			wantErr:   false,
			wantMatch: true,
		},
		{
			name: "This config does not generate the example/versioned/testing/fields.go file",
			fields: fields{
				StructName: StructName{
					Lower: "testing",
				},
				OutputDir: tempDir,
				ProcessedFields: []FieldInfo{
					{
						Name:          "InEveryVersion",
						FormattedName: "InEveryVersion",
						Type:          "*string",
						Tag:           "json:\"in_every_version\"",
					},
				},
			},
			wantErr:   false,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				StructName:      tt.fields.StructName,
				OutputDir:       tt.fields.OutputDir,
				ProcessedFields: tt.fields.ProcessedFields,
			}
			if err := g.FieldsFile(); (err != nil) != tt.wantErr {
				t.Errorf("FieldsFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Path to the generated file
			generatedFilePath := filepath.Join(tempDir, "versioned", tt.fields.StructName.Lower, "fields.go")

			// Read the contents of the generated file
			generatedFileContent, err := os.ReadFile(generatedFilePath)
			assert.NoError(t, err)

			// Path to the reference file
			referenceFilePath := filepath.Join("example", "versioned", tt.fields.StructName.Lower, "fields.go")

			// Read the contents of the reference file
			referenceFileContent, err := os.ReadFile(referenceFilePath)
			assert.NoError(t, err)

			// Compare the contents
			if tt.wantMatch {
				assert.Equal(t, string(referenceFileContent), string(generatedFileContent), "The generated file content does not match the reference file content")
			} else {
				assert.NotEqual(t, string(referenceFileContent), string(generatedFileContent), "The generated file content should not match the reference file content")
			}
		})
	}
}
