package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	VersionTag = "version"
)

type Format struct {
	Versions       map[int][]string
	SortedVersions []int
	CustomType     bool
}

func (f *Format) FieldType(expr ast.Expr, pointer bool) string {
	result := ""
	if pointer {
		result += "*"
	}

	switch t := expr.(type) {
	case *ast.Ident:
		// Check if it's a custom type (non-builtin)
		if isCustomType(t.Name) {
			f.CustomType = true
			result += "originalPackage." + t.Name
		} else {
			result += t.Name
		}
	case *ast.ArrayType:
		// For arrays and slices
		elementType := f.FieldType(t.Elt, false)
		result += fmt.Sprintf("[]%s", elementType)
	case *ast.MapType:
		// For maps
		keyType := f.FieldType(t.Key, false)
		valueType := f.FieldType(t.Value, false)
		result += fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.StarExpr:
		// For pointers, simply add an extra star
		pointedType := f.FieldType(t.X, true)
		result += pointedType
	case *ast.SelectorExpr:
		// For qualified identifiers (e.g., time.Time)
		result += fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	default:
		// Fallback for other types
		result += "any" // Generic pointer type as fallback
	}

	return result
}

func isCustomType(typeName string) bool {
	builtInTypes := map[string]bool{
		"bool": true, "int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"uintptr": true, "float32": true, "float64": true, "complex64": true, "complex128": true,
		"string": true, "byte": true, "rune": true, "error": true, "any": true, "interface{}": true,
	}

	// keep trimming the pointer until it's not a pointer anymore
	for strings.HasPrefix(typeName, "*") {
		typeName = typeName[1:]
	}
	_, isBuiltIn := builtInTypes[typeName]
	return !isBuiltIn
}

func (f *Format) IdentifyVersions(structType *ast.StructType) {
	var allTags []string
	// Collect all version tags from the struct fields
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get(VersionTag)
			allTags = append(allTags, tag)
		}
	}

	maxVersion := f.DetermineMaxVersion(allTags)

	versionMap := make(map[int][]string)
	for _, field := range structType.Fields.List {
		var versions []int
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get(VersionTag)
			versions = f.ParseVersionTag(tag, maxVersion)
		} else {
			for v := 1; v <= maxVersion; v++ {
				versions = append(versions, v)
			}
		}

		fieldType := fmt.Sprintf("%s", field.Type)
		fieldType = f.FieldType(field.Type, false)
		for _, name := range field.Names {
			fieldStr := fmt.Sprintf("%s %s", name, fieldType)
			for _, v := range versions {
				versionMap[v] = append(versionMap[v], fieldStr)
			}
		}
	}

	f.Versions = versionMap

	for versionItem := range versionMap {
		f.SortedVersions = append(f.SortedVersions, versionItem)
	}
	sort.Ints(f.SortedVersions)
}

func (f *Format) ParseVersionTag(tag string, maxVersion int) []int {
	if tag == "" {
		// If no tag, include in all versions
		var versions []int
		for i := 1; i <= maxVersion; i++ {
			versions = append(versions, i)
		}
		return versions
	}

	// Remaining logic stays the same
	start, end, err := f.ParseVersionRange(tag)
	if err != nil || start > maxVersion {
		return []int{}
	}

	if end == -1 { // No upper limit specified
		end = maxVersion
	}

	var versions []int
	for v := start; v <= end; v++ {
		versions = append(versions, v)
	}

	return versions
}

func (f *Format) DetermineMaxVersion(versionTags []string) int {
	maxVersion := 1

	for _, tag := range versionTags {
		_, end, err := f.ParseVersionRange(tag)
		if err == nil && end > maxVersion {
			maxVersion = end
		}
	}

	return maxVersion
}

func (f *Format) ParseVersionRange(tag string) (int, int, error) {
	if tag == "" {
		return 1, 1, nil // Default to version 1 if no tag
	}

	if strings.Contains(tag, "+") {
		// For "2+" style tags
		start, err := strconv.Atoi(strings.TrimSuffix(tag, "+"))
		if err != nil {
			return 0, 0, err // Error in parsing the tag
		}
		return start, -1, nil // -1 indicates no upper limit
	}

	if strings.Contains(tag, "-") {
		// For "1-3" or "-3" style tags
		parts := strings.Split(tag, "-")
		start, end := 1, 0
		var err error

		if parts[0] != "" {
			start, err = strconv.Atoi(parts[0])
			if err != nil {
				return 0, 0, err // Error in parsing the tag
			}
		}

		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err // Error in parsing the tag
		}

		return start, end, nil
	}

	// For single version tags like "2"
	version, err := strconv.Atoi(tag)
	if err != nil {
		return 0, 0, err // Error in parsing the tag
	}
	return version, version, nil
}

func (f *Format) ExcludeVersionTag(tag string) string {
	var result []string
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		if !strings.HasPrefix(t, fmt.Sprintf("%s:", VersionTag)) {
			result = append(result, t)
		}
	}
	return strings.Join(result, " ")
}
