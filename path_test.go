package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestResolver_FindGoModPath_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		assert.NoError(t, err)
	}(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte("module example.com/my/module\n"), 0644)
	assert.NoError(t, err)

	resolver := Resolver{}
	err = resolver.FindGoModPath(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, goModPath, resolver.GoModPath)
}

func TestResolver_FindGoModPath_Failure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		assert.NoError(t, err)
	}(tempDir)

	resolver := Resolver{}
	err = resolver.FindGoModPath(tempDir)
	assert.Error(t, err)
}

func TestResolver_GetBaseImportPath_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		assert.NoError(t, err)
	}(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte("module example.com/my/module\n"), 0644)
	assert.NoError(t, err)

	resolver := Resolver{GoModPath: goModPath}
	err = resolver.GetBaseImportPath()
	assert.NoError(t, err)
	assert.Equal(t, "example.com/my/module", resolver.ImportPath)
}

func TestResolver_GetBaseImportPath_Failure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		assert.NoError(t, err)
	}(tempDir)

	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte("// This is a go.mod file without a module declaration\n"), 0644)
	assert.NoError(t, err)

	resolver := Resolver{GoModPath: goModPath}
	err = resolver.GetBaseImportPath()
	assert.Error(t, err)
}
