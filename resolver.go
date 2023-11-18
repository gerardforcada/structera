package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Resolver struct {
	ImportPath string
	GoModPath  string
}

func (r *Resolver) GetBaseImportPath(goModFilePath string) error {
	data, err := os.ReadFile(goModFilePath)
	if err != nil {
		return err
	}

	modulePrefix := "module "
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, modulePrefix) {
			r.ImportPath = strings.TrimSpace(line[len(modulePrefix):])
			return nil
		}
	}

	return fmt.Errorf("module declaration not found in go.mod")
}

func (r *Resolver) FindGoModPath(startDir string) error {
	currentDir := startDir
	for {
		files, err := os.ReadDir(currentDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if f.Name() == "go.mod" {
				r.GoModPath = filepath.Join(currentDir, f.Name())
				return nil
			}
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return fmt.Errorf("go.mod not found")
}
