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
			expected: "NestedStruct",
		},
		{
			name:     "Pointer to Struct Type",
			expr:     &ast.StructType{},
			pointer:  true,
			expected: "*NestedStruct",
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
			pointer:  false,
			expected: "*int",
		},
		{
			name: "Pointer to Pointer",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "*int"},
			},
			pointer:  true,
			expected: "**int",
		},
		{
			name: "Double Pointer",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  false,
			expected: "**string",
		},
		{
			name: "Pointer to Double Pointer",
			expr: &ast.StarExpr{
				X: &ast.StarExpr{
					X: &ast.Ident{Name: "**string"},
				},
			},
			pointer:  true,
			expected: "***string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFieldType(tt.expr, tt.pointer)
			assert.Equal(t, tt.expected, result)
		})
	}
}
