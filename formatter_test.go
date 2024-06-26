package main

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFieldType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		pointer  bool
		expected string
	}{
		{
			name:     "Simple Type",
			expr:     &ast.Ident{Name: "int"},
			pointer:  false,
			expected: "int",
		},
		{
			name:     "Pointer to Simple Type",
			expr:     &ast.Ident{Name: "string"},
			pointer:  true,
			expected: "*string",
		},
		{
			name: "Array Type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "byte"},
			},
			pointer:  false,
			expected: "[]byte",
		},
		{
			name: "Pointer to Array Type",
			expr: &ast.ArrayType{
				Elt: &ast.Ident{Name: "byte"},
			},
			pointer:  true,
			expected: "*[]byte",
		},
		{
			name: "Map Type",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			pointer:  false,
			expected: "map[string]int",
		},
		{
			name: "Pointer to Map Type",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "int"},
				Value: &ast.Ident{Name: "string"},
			},
			pointer:  true,
			expected: "*map[int]string",
		},
		{
			name:     "Nested Struct Type",
			expr:     &ast.StructType{},
			pointer:  false,
			expected: "any",
		},
		{
			name:     "Pointer to Struct Type",
			expr:     &ast.StructType{},
			pointer:  true,
			expected: "*any",
		},
		{
			name: "Qualified Identifier",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			pointer:  false,
			expected: "time.Time",
		},
		{
			name: "Pointer to Qualified Identifier",
			expr: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "time"},
				Sel: &ast.Ident{Name: "Time"},
			},
			pointer:  true,
			expected: "*time.Time",
		},
		{
			name: "Pointer",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "*int"},
			},
			pointer:  true,
			expected: "***int",
		},
		{
			name: "Pointer2",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "*int"},
			},
			pointer:  false,
			expected: "**int",
		},
		{
			name: "Double Pointer",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  true,
			expected: "*****string",
		},
		{
			name: "Double Pointer2",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  false,
			expected: "****string",
		},
		{
			name: "Triple Pointer",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  true,
			expected: "*****string",
		},
		{
			name: "Triple Pointer2",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  false,
			expected: "****string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Format{}
			result := v.FieldType(tt.expr, tt.pointer)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_IdentifyVersions(t *testing.T) {
	tests := []struct {
		name     string
		expected map[int][]string
	}{
		{
			name: "Basic test",
			expected: map[int][]string{
				1: {"InEveryVersion string", "OnlyIn1 int", "FromStartTo3 []byte", "From1to4 float32"},
				2: {"InEveryVersion string", "From2ToEnd uint8", "FromStartTo3 []byte", "From1to4 float32"},
				3: {"InEveryVersion string", "From2ToEnd uint8", "FromStartTo3 []byte", "From1to4 float32"},
				4: {"InEveryVersion string", "From2ToEnd uint8", "From1to4 float32"},
				5: {"InEveryVersion string", "From2ToEnd uint8", "OnlyIn5 int32"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Format{}
			structType := &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "InEveryVersion"}},
							Type:  &ast.Ident{Name: "string"},
							Tag:   &ast.BasicLit{Value: "`json:\"in_every_version\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "OnlyIn1"}},
							Type:  &ast.Ident{Name: "int"},
							Tag:   &ast.BasicLit{Value: "`version:\"1\" json:\"only_in_1\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "From2ToEnd"}},
							Type:  &ast.Ident{Name: "uint8"},
							Tag:   &ast.BasicLit{Value: "`version:\"2+\" json:\"from_2_to_end\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "FromStartTo3"}},
							Type:  &ast.ArrayType{Elt: &ast.Ident{Name: "byte"}},
							Tag:   &ast.BasicLit{Value: "`version:\"-3\" json:\"from_start_to_3\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "From1to4"}},
							Type:  &ast.Ident{Name: "float32"},
							Tag:   &ast.BasicLit{Value: "`version:\"1-4\" json:\"from_1_to_4\"`"},
						},
						{
							Names: []*ast.Ident{{Name: "OnlyIn5"}},
							Type:  &ast.Ident{Name: "int32"},
							Tag:   &ast.BasicLit{Value: "`version:\"5\" json:\"only_in_5\"`"},
						},
					},
				},
			}
			v.IdentifyVersions(structType)
			assert.Equal(t, tt.expected, v.Versions)
		})
	}
}

func TestVersion_ParseVersionTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		maxVersion int
		expected   []int
	}{
		{"Empty tag", "", 3, []int{1, 2, 3}},
		{"Single version", "2", 3, []int{2}},
		{"Range", "1-3", 5, []int{1, 2, 3}},
		{"Open-ended", "2+", 3, []int{2, 3}},
		{"Invalid tag", "abc", 3, []int{}},
	}

	v := Format{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ParseVersionTag(tt.tag, tt.maxVersion)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersion_DetermineMaxVersion(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected int
	}{
		{"Single version", []string{"1"}, 1},
		{"Range", []string{"1-3"}, 3},
		{"Multiple tags", []string{"1", "2", "3-4"}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Format{}
			maxVersion := v.DetermineMaxVersion(tt.tags)
			assert.Equal(t, tt.expected, maxVersion)
		})
	}
}

func TestVersion_ParseVersionRange(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "Single version",
			tag:       "2",
			wantStart: 2,
			wantEnd:   2,
			wantErr:   false,
		},
		{
			name:      "Version range",
			tag:       "1-3",
			wantStart: 1,
			wantEnd:   3,
			wantErr:   false,
		},
		{
			name:      "Open-ended range",
			tag:       "4+",
			wantStart: 4,
			wantEnd:   -1,
			wantErr:   false,
		},
		{
			name:      "Invalid format",
			tag:       "abc",
			wantStart: 0,
			wantEnd:   0,
			wantErr:   true,
		},
		{
			name:      "Empty tag",
			tag:       "",
			wantStart: 1,
			wantEnd:   1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Format{}
			gotStart, gotEnd, err := v.ParseVersionRange(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStart != tt.wantStart {
				t.Errorf("ParseVersionRange() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("ParseVersionRange() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestVersion_ExcludeVersionTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{"Single version tag", `version:"1" json:"field1"`, `json:"field1"`},
		{"Multiple tags", `json:"field1" version:"1-2" xml:"field1"`, `json:"field1" xml:"field1"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Format{}
			result := v.ExcludeVersionTag(tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}
