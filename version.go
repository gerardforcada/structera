package main

import (
	"fmt"
	"github.com/gerardforcada/structera/helpers"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

type Version struct {
	Versions map[int][]string
}

func (v *Version) IdentifyVersions(structType *ast.StructType) {
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

	v.Versions = versionMap
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

func (v *Version) ExcludeVersionTag(tag string) string {
	var result []string
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		if !strings.HasPrefix(t, "version:") {
			result = append(result, t)
		}
	}
	return strings.Join(result, " ")
}
