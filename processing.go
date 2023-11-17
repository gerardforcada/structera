package main

import (
	"fmt"
	"github.com/gerardforcada/structera/helpers"
	"go/ast"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func identifyVersions(structType *ast.StructType) map[int][]string {
	var allTags []string
	// Collect all version tags from the struct fields
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("version")
			allTags = append(allTags, tag)
		}
	}

	maxVersion := determineMaxVersion(allTags)

	versionMap := make(map[int][]string)
	for _, field := range structType.Fields.List {
		var versions []int
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("version")
			versions = parseVersionTag(tag, maxVersion)
		} else {
			for v := 1; v <= maxVersion; v++ {
				versions = append(versions, v)
			}
		}

		fieldType := fmt.Sprintf("%s", field.Type)
		fieldType = formatFieldType(field.Type, false)
		for _, name := range field.Names {
			fieldStr := fmt.Sprintf("%s %s", name, fieldType)
			for _, v := range versions {
				versionMap[v] = append(versionMap[v], fieldStr)
			}
		}
	}

	return versionMap
}

func parseVersionTag(tag string, maxVersion int) []int {
	var versions []int

	if tag == "" {
		return []int{1} // Default to version 1 if no tag
	}

	if strings.Contains(tag, "+") {
		// For "2+" style tags
		minVersion, err := strconv.Atoi(strings.TrimSuffix(tag, "+"))
		if err != nil {
			return []int{} // Error in parsing the tag
		}
		for v := minVersion; v <= maxVersion; v++ {
			versions = append(versions, v)
		}
	} else if strings.Contains(tag, "-") {
		// For "1-3" or "-3" style tags
		parts := strings.Split(tag, "-")
		var start, end int
		var err error
		if parts[0] == "" { // For "-3" style tags
			start = 1
			end, err = strconv.Atoi(parts[1])
		} else { // For "1-3" style tags
			start, err = strconv.Atoi(parts[0])
			end, err = strconv.Atoi(parts[1])
		}
		if err != nil {
			return []int{} // Error in parsing the tag
		}
		for v := start; v <= end; v++ {
			versions = append(versions, v)
		}
	} else {
		// For single version tags like "2"
		version, err := strconv.Atoi(tag)
		if err != nil {
			return []int{} // Error in parsing the tag
		}
		versions = []int{version}
	}

	return versions
}

func determineMaxVersion(versionTags []string) int {
	maxVersion := 1 // Default to version 1 if no higher versions are found

	for _, tag := range versionTags {
		if strings.Contains(tag, "+") {
			// For "2+" style tags, assume a large number since we don't have the upper limit
			parts := strings.Split(tag, "+")
			value, err := strconv.Atoi(parts[0])
			if err != nil {
				continue // Error in parsing the tag
			}
			maxVersion = helpers.Max(maxVersion, value)
		} else if strings.Contains(tag, "-") {
			// For "1-3" or "-3" style tags
			parts := strings.Split(tag, "-")
			var end int
			var err error
			if parts[0] == "" { // For "-3" style tags
				end, err = strconv.Atoi(parts[1])
			} else { // For "1-3" style tags
				end, err = strconv.Atoi(parts[1])
			}
			if err == nil {
				maxVersion = helpers.Max(maxVersion, end)
			}
		} else {
			// For single version tags like "2"
			v, err := strconv.Atoi(tag)
			if err == nil {
				maxVersion = helpers.Max(maxVersion, v)
			}
		}
	}

	return maxVersion
}

func processStruct(structName string, structType *ast.StructType, existingImports []string, demoImportPath string) (string, error) {
	versions := identifyVersions(structType)
	if len(versions) == 0 {
		return "", fmt.Errorf("no version tags found in struct")
	}

	// Sort the version numbers
	var versionNumbers []int
	for v := range versions {
		versionNumbers = append(versionNumbers, v)
	}
	sort.Ints(versionNumbers)

	var buf strings.Builder
	buf.WriteString("package versioned\n\n")
	buf.WriteString("import (\n")
	buf.WriteString(fmt.Sprintf("\t\"%s/version\"\n", LibraryPackage))
	buf.WriteString(fmt.Sprintf("\t\"%s\"\n", demoImportPath)) // Add dynamic import path for demo

	// Include existing imports from the original file
	for _, imp := range existingImports {
		buf.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}
	buf.WriteString(")\n\n")

	// DemoVersions struct
	buf.WriteString("type DemoVersions struct {\n")
	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\tV%d demo.V%d[DemoV%d]\n", v, v, v))
	}
	buf.WriteString("}\n\n")

	// Demo struct
	buf.WriteString("type Demo struct {\n\tdemo.PointerFields\n\tDemoVersions\n}\n\n")

	// Initialize function
	buf.WriteString(fmt.Sprintf("func (d *%s) Initialize() {\n", structName))
	buf.WriteString("\td.DemoVersions = DemoVersions{\n")

	for _, v := range versionNumbers {
		buf.WriteString(fmt.Sprintf("\t\tV%d: &%sV%d{},\n", v, structName, v))
	}

	buf.WriteString("\t}\n}\n\n")

	// Additional methods for Demo
	buf.WriteString(generateDemoMethods(structName, versionNumbers))

	// Version-specific struct types and methods
	for _, v := range versionNumbers {
		fields, ok := versions[v]
		if !ok {
			continue // Skip if version number is not found in the map
		}

		// Extract tags from original struct fields
		tags := make(map[string]string)
		for _, field := range structType.Fields.List {
			if field.Tag != nil {
				for _, name := range field.Names {
					tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json")
					if tag != "" {
						tags[name.Name] = fmt.Sprintf("json:\"%s\"", tag)
					}
				}
			}
		}

		buf.WriteString(fmt.Sprintf("type DemoV%d struct {\n", v))
		buf.WriteString(formatStructFields(fields, tags))
		buf.WriteString("}\n\n")
		buf.WriteString(generateVersionMethods(v))
	}

	return buf.String(), nil
}
