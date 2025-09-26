package scaffolder

import (
	"grei-cli/internal/adapters/filesystem"
	"grei-cli/internal/core/recipe"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffold_GoCli(t *testing.T) {
	// Arrange
	fsMock := filesystem.NewMockRepository()
	defer fsMock.Clean()

	fsMock.AddManifest("code", "go-cli", "Go", "Cobra", "")
	fsMock.AddTemplate("templates/generic/README.md.tmpl", "README for {{ .Project.Name }}")
	fsMock.AddTemplate("templates/generic/.gitignore.tmpl", "*.log")
	fsMock.AddTemplate("templates/go-cli/Makefile.tmpl", "BINARY_NAME={{ .Project.Name }}")

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestCliProject",
			Type: "cli",
		},
		Stack: recipe.Stack{
			Language: "Go",
			Tooling:  "Cobra",
		},
	}

	service := NewService(fsMock)
	tmpDir := fsMock.TempDir()

	// Act
	err := service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err != nil {
		t.Fatalf("Scaffold() returned an unexpected error: %v", err)
	}

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
	fsMock := filesystem.NewMockRepository()
	defer fsMock.Clean()

	fsMock.AddManifest("persistence", "postgresql", "", "", "postgresql")
	fsMock.AddTemplate("templates/generic/README.md.tmpl", "README for {{ .Project.Name }}")
	fsMock.AddTemplate("templates/postgresql/docker-compose.yml.tmpl", `
services:
  db:
    image: postgres
    environment:
      POSTGRES_DB: {{ .Project.Name | ToLower }}_db
      POSTGRES_USER: {{ .Project.Name | ToLower }}_user
      POSTGRES_PASSWORD: {{ .Project.Name | ToLower }}_password
`)

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestPostgresProject",
		},
		Persistence: recipe.Persistence{
			Type: "postgresql",
		},
	}

	service := NewService(fsMock)
	tmpDir := fsMock.TempDir()

	// Act
	err := service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err != nil {
		t.Fatalf("Scaffold() returned an unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "docker-compose.yml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it was not: %s", path)
	}

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
	fsMock := filesystem.NewMockRepository()
	defer fsMock.Clean()

	fsMock.AddTemplate("templates/generic/README.md.tmpl", "README for {{ .Project.Name }}")
	fsMock.AddTemplate("templates/generic/.gitignore.tmpl", "*.log")

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestGenericProject",
			Type: "generic",
		},
	}

	service := NewService(fsMock)
	tmpDir := fsMock.TempDir()

	// Act
	err := service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err != nil {
		t.Fatalf("Scaffold() returned an unexpected error: %v", err)
	}

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
	fsMock := filesystem.NewMockRepository()
	defer fsMock.Clean()

	fsMock.AddTemplate("templates/generic/README.md.tmpl", "README for {{ .Project.Name }}")

	projRecipe := &recipe.Recipe{
		Project: recipe.Project{
			Name: "TestWriteErrorProject",
			Type: "generic",
		},
	}

	service := NewService(fsMock)
	tmpDir := fsMock.TempDir()

	// Make the directory read-only to cause a write error.
	if err := os.Chmod(tmpDir, 0555); err != nil {
		t.Fatalf("Failed to change permissions of temp dir: %v", err)
	}

	// Act
	err := service.Scaffold(tmpDir, projRecipe)

	// Assert
	if err == nil {
		t.Errorf("Scaffold() should have returned an error, but it did not")
	}
}
