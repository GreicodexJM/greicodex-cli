package scaffolder

import (
	"grei-cli/internal/core/recipe"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffold_GoCli(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestCliProject",
			Type: "go-cli",
		},
	}

	service := NewService()

	// Act
	err = service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err != nil {
		t.Fatalf("Scaffold() returned an unexpected error: %v", err)
	}

	// Check that the expected files were created.
	expectedFiles := []string{
		"README.md",
		".gitignore",
		"Makefile",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file to be created, but it was not: %s", path)
		}
	}

	// Check that the Makefile content was correctly templated.
	makefileContent, err := os.ReadFile(filepath.Join(tmpDir, "Makefile"))
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}

	if !strings.Contains(string(makefileContent), "BINARY_NAME=TestCliProject") {
		t.Errorf("Expected Makefile to contain the correct binary name, but it did not.")
	}
}
