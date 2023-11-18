package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

const (
	VersionTag = "version"
)

type Version struct {
	Versions map[int][]string
}

func (v *Version) IdentifyVersions(structType *ast.StructType) {
	var allTags []string
	// Collect all version tags from the struct fields
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get(VersionTag)
			allTags = append(allTags, tag)
		}
	}

	maxVersion := v.determineMaxVersion(allTags)

	versionMap := make(map[int][]string)
	for _, field := range structType.Fields.List {
		var versions []int
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get(VersionTag)
			versions = v.parseVersionTag(tag, maxVersion)
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

func (v *Version) parseVersionTag(tag string, maxVersion int) []int {
	start, end, err := v.parseVersionRange(tag)
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

func (v *Version) determineMaxVersion(versionTags []string) int {
	maxVersion := 1

	for _, tag := range versionTags {
		_, end, err := v.parseVersionRange(tag)
		if err == nil && end > maxVersion {
			maxVersion = end
		}
	}

	return maxVersion
}

func (v *Version) parseVersionRange(tag string) (int, int, error) {
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

func (v *Version) ExcludeVersionTag(tag string) string {
	var result []string
	tags := strings.Split(tag, " ")
	for _, t := range tags {
		if !strings.HasPrefix(t, fmt.Sprintf("%s:", VersionTag)) {
			result = append(result, t)
		}
	}
	return strings.Join(result, " ")
}
