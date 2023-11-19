package main

import (
	"fmt"
	"go/ast"
)

func formatFieldType(expr ast.Expr, pointer bool) string {
	result := ""
	if pointer {
		result += "*"
	}

	switch t := expr.(type) {
	case *ast.Ident:
		result += t.Name // Simple types like int, string, etc.
	case *ast.ArrayType:
		// For arrays and slices
		elementType := formatFieldType(t.Elt, false) // No extra pointer for element type
		result += fmt.Sprintf("[]%s", elementType)
	case *ast.MapType:
		// For maps
		keyType := formatFieldType(t.Key, false)     // No extra pointer for key type
		valueType := formatFieldType(t.Value, false) // No extra pointer for value type
		result += fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.StarExpr:
		// For pointers, simply add an extra star
		pointedType := formatFieldType(t.X, false)
		result += fmt.Sprintf("%s", pointedType)
	case *ast.StructType:
		// Handle nested structs
		result += "NestedStruct" // Adjust as needed for nested structs
	case *ast.SelectorExpr:
		// For qualified identifiers (e.g., time.Time)
		result += fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	default:
		// Fallback for other types
		result += "interface{}" // Generic pointer type as fallback
	}

	return result
}
