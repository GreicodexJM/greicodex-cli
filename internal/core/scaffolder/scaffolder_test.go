package scaffolder

import (
	"grei-cli/internal/core/recipe"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffold_GoCli(t *testing.T) {
	// Debugging: List all embedded files.
	fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		t.Logf("Found embedded file: %s", path)
		return nil
	})

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

func TestScaffold_Postgresql(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestPostgresProject",
		},
		Persistence: recipe.Persistence{
			Type: "postgresql",
		},
	}

	service := NewService()

	// Act
	err = service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err != nil {
		t.Fatalf("Scaffold() returned an unexpected error: %v", err)
	}

	// Check that the docker-compose file was created.
	path := filepath.Join(tmpDir, "docker-compose.yml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it was not: %s", path)
	}

	// Check that the content was correctly templated.
	composeContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}

	if !strings.Contains(string(composeContent), "POSTGRES_DB: testpostgresproject_db") {
		t.Errorf("Expected docker-compose.yml to contain the correct db name, but it did not.")
	}
	if !strings.Contains(string(composeContent), "POSTGRES_USER: testpostgresproject_user") {
		t.Errorf("Expected docker-compose.yml to contain the correct user, but it did not.")
	}
	if !strings.Contains(string(composeContent), "POSTGRES_PASSWORD: testpostgresproject_password") {
		t.Errorf("Expected docker-compose.yml to contain the correct password, but it did not.")
	}
}

func TestScaffold_Generic(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestGenericProject",
			Type: "generic",
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
	}

	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file to be created, but it was not: %s", path)
		}
	}
}

func TestScaffold_WriteError(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Make the directory read-only to cause a write error.
	if err := os.Chmod(tmpDir, 0555); err != nil {
		t.Fatalf("Failed to change permissions of temp dir: %v", err)
	}

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestWriteErrorProject",
			Type: "generic",
		},
	}

	service := NewService()

	// Act
	err = service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err == nil {
		t.Errorf("Scaffold() should have returned an error, but it did not")
	}
}

func TestScaffold_GoCli_WriteError(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Make the directory read-only to cause a write error.
	if err := os.Chmod(tmpDir, 0555); err != nil {
		t.Fatalf("Failed to change permissions of temp dir: %v", err)
	}

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestGoCliWriteErrorProject",
			Type: "go-cli",
		},
	}

	service := NewService()

	// Act
	err = service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err == nil {
		t.Errorf("Scaffold() should have returned an error, but it did not")
	}
}

func TestScaffold_Postgresql_WriteError(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Make the directory read-only to cause a write error.
	if err := os.Chmod(tmpDir, 0555); err != nil {
		t.Fatalf("Failed to change permissions of temp dir: %v", err)
	}

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestPostgresqlWriteErrorProject",
		},
		Persistence: recipe.Persistence{
			Type: "postgresql",
		},
	}

	service := NewService()

	// Act
	err = service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err == nil {
		t.Errorf("Scaffold() should have returned an error, but it did not")
	}
}
