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

func (r *Resolver) FindGoModPath(startDir string) error {
	currentDir := startDir
	for {
		// Check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			r.GoModPath = filepath.Join(currentDir, "go.mod")
			return nil
		}

		// Move up to the parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return fmt.Errorf("go.mod not found")
		}
		currentDir = parentDir
	}
}

func (r *Resolver) GetBaseImportPath() error {
	data, err := os.ReadFile(r.GoModPath)
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
