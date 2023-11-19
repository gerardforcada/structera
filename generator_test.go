package main

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_FileFromTemplate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)

	// Clean up after the test
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	// Test case
	test := struct {
		name    string
		input   GenerateFileFromTemplateInput
		wantErr bool
	}{
		name: "This config generates the example/versioned/testing/fields.go file",
		input: GenerateFileFromTemplateInput{
			TemplateFilePath: "fields.go.tmpl",
			OutputFilePath:   filepath.Join(tempDir, "fields.go"),
			Data: FieldsTemplateData{
				PackageName: "testing",
				Fields: []FieldInfo{
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
		},
		wantErr: false,
	}

	t.Run(test.name, func(t *testing.T) {
		g := &Generator{} // Assuming Generator doesn't need additional setup for this method

		if err := g.FileFromTemplate(test.input); (err != nil) != test.wantErr {
			t.Errorf("FileFromTemplate() error = %v, wantErr %v", err, test.wantErr)
		}

		// Verify the file was created
		_, err := os.Stat(test.input.OutputFilePath)
		assert.NoError(t, err, "The file was not created")

		// Read the contents of the generated file
		generatedFileContent, err := os.ReadFile(test.input.OutputFilePath)
		assert.NoError(t, err)

		// Path to the reference file
		referenceFilePath := filepath.Join("example", "versioned", "testing", "fields.go")

		// Read the contents of the reference file
		referenceFileContent, err := os.ReadFile(referenceFilePath)
		assert.NoError(t, err)

		// Compare the contents
		assert.Equal(t, string(referenceFileContent), string(generatedFileContent), "The generated file content does not match the reference file content")
	})
}

func TestGenerator_PrepareVersionedFields(t *testing.T) {
	tests := []struct {
		name            string
		version         *Version
		processedFields []FieldInfo
		want            map[int][]FieldInfo
	}{
		{
			name: "Basic Test",
			version: &Version{
				Versions: map[int][]string{
					1: {"Field1 string", "Field2 int"},
					2: {"Field2 int", "Field3 float64"},
				},
			},
			processedFields: []FieldInfo{
				{Name: "Field1", Type: "*string"},
				{Name: "Field2", Type: "*int"},
				{Name: "Field3", Type: "*float64"},
			},
			want: map[int][]FieldInfo{
				1: {{Name: "Field1", Type: "string"}, {Name: "Field2", Type: "int"}},
				2: {{Name: "Field2", Type: "int"}, {Name: "Field3", Type: "float64"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Version:         tt.version,
				ProcessedFields: tt.processedFields,
			}

			g.PrepareVersionedFields()

			assert.Equal(t, tt.want, g.VersionedFields)
		})
	}
}

func TestGenerator_ProcessFieldInfo(t *testing.T) {
	tests := []struct {
		name           string
		structType     *ast.StructType
		version        *Version
		expectedFields []FieldInfo
		expectedMaxLen int
	}{
		{
			name: "Basic Struct",
			structType: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Field1"}},
							Type:  &ast.Ident{Name: "string"},
							Tag:   &ast.BasicLit{Value: "`json:\"field1\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "Field2"}},
							Type:  &ast.Ident{Name: "int"},
						},
						{
							Names: []*ast.Ident{{Name: "Field3"}},
							Type:  &ast.Ident{Name: "int"},
							Tag:   &ast.BasicLit{Value: "`version:\"1\"`"},
						},
					},
				},
			},
			expectedFields: []FieldInfo{
				{Name: "Field1", Type: "*string", Tag: "json:\"field1\""},
				{Name: "Field2", Type: "*int"},
				{Name: "Field3", Type: "*int"},
			},
			expectedMaxLen: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Version: tt.version,
			}

			fields, maxLen, err := g.ProcessFieldInfo(tt.structType)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFields, fields)
			assert.Equal(t, tt.expectedMaxLen, maxLen)
		})
	}
}
