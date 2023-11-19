package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_VersionFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// Clean up after the test
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	type fields struct {
		StructName StructName
		OutputDir  string
	}
	tests := []struct {
		name      string
		fields    fields
		versions  []int
		wantErr   bool
		wantMatch bool
	}{
		{
			name: "This config generates the example/versioned/testing/version.go file",
			fields: fields{
				StructName: StructName{
					Lower: "testing",
				},
				OutputDir: tempDir,
			},
			versions:  []int{1, 2, 3, 4},
			wantErr:   false,
			wantMatch: true,
		},
		{
			name: "This config does not generate the example/versioned/testing/version.go file",
			fields: fields{
				StructName: StructName{
					Lower: "testing",
				},
				OutputDir: tempDir,
			},
			versions:  []int{1},
			wantErr:   false,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				StructName: tt.fields.StructName,
				OutputDir:  tt.fields.OutputDir,
			}
			if err := g.VersionFile(tt.versions); (err != nil) != tt.wantErr {
				t.Errorf("VersionFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Path to the generated file
			generatedFilePath := filepath.Join(tempDir, "versioned", tt.fields.StructName.Lower, "version.go")

			// Read the contents of the generated file
			generatedFileContent, err := os.ReadFile(generatedFilePath)
			assert.NoError(t, err)

			// Path to the reference file
			referenceFilePath := filepath.Join("example", "versioned", tt.fields.StructName.Lower, "version.go")

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
