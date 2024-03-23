package main

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestGenerator_PrepareVersionedFields(t *testing.T) {
	tests := []struct {
		name            string
		format          *Format
		processedFields []HubFieldInfo
		want            map[int][]HubFieldInfo
	}{
		{
			name: "Basic Test",
			format: &Format{
				Versions: map[int][]string{
					1: {"Field1 string", "Field2 int"},
					2: {"Field2 int", "Field3 float64"},
				},
			},
			processedFields: []HubFieldInfo{
				{Name: "Field1", Type: "*string"},
				{Name: "Field2", Type: "*int"},
				{Name: "Field3", Type: "*float64"},
			},
			want: map[int][]HubFieldInfo{
				1: {{Name: "Field1", Type: "string"}, {Name: "Field2", Type: "int"}},
				2: {{Name: "Field2", Type: "int"}, {Name: "Field3", Type: "float64"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Format:          tt.format,
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
		format         *Format
		expectedFields []HubFieldInfo
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
			expectedFields: []HubFieldInfo{
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
				Format: tt.format,
			}

			fields, maxLen, err := g.ProcessFieldInfo(tt.structType)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFields, fields)
			assert.Equal(t, tt.expectedMaxLen, maxLen)
		})
	}
}
