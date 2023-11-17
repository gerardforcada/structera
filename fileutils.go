package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getBaseImportPath(goModFilePath string) (string, error) {
	data, err := os.ReadFile(goModFilePath)
	if err != nil {
		return "", err
	}

	modulePrefix := "module "
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, modulePrefix) {
			return strings.TrimSpace(line[len(modulePrefix):]), nil
		}
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}

func findGoModPath(startDir string) (string, error) {
	currentDir := startDir
	for {
		files, err := os.ReadDir(currentDir)
		if err != nil {
			return "", err
		}

		for _, f := range files {
			if f.Name() == "go.mod" {
				return filepath.Join(currentDir, f.Name()), nil
			}
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("go.mod not found")
}
