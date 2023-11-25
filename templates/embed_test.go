package templates

import (
	"testing"
)

// TestFS checks if the embedded file system can be accessed and specific files exist.
func TestFS(t *testing.T) {
	expectedFiles := []string{"struct.go.tmpl", "types.go.tmpl"}
	notExpectedFiles := []string{"embed_test.go"}

	for _, fileName := range expectedFiles {
		_, err := FS.ReadFile(fileName)
		if err != nil {
			t.Errorf("Failed to read %s from embedded file system: %v", fileName, err)
		}
	}

	for _, fileName := range notExpectedFiles {
		_, err := FS.ReadFile(fileName)
		if err == nil {
			t.Errorf("Expected to fail reading %s from embedded file system, but it did not", fileName)
		}
	}
}
