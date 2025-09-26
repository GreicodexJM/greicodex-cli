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
	fsMock.AddTemplate("skeletons/generic/README.md", "README for {{ .Project.Name }}")
	fsMock.AddTemplate("skeletons/generic/.gitignore", "*.log")
	fsMock.AddTemplate("skeletons/go-cli/Makefile", "BINARY_NAME={{ .Project.Name }}")
	fsMock.AddTemplate("skeletons/generic/manifest.json", `{"files": [{"path": "README.md", "strategy": "overwrite"}, {"path": ".gitignore", "strategy": "overwrite"}]}`)
	fsMock.AddTemplate("skeletons/go-cli/manifest.json", `{"files": [{"path": "Makefile", "strategy": "overwrite"}]}`)

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
	err := service.Scaffold(tmpDir, fsMock.TempDir(), projRecipe)

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
	fsMock.AddTemplate("skeletons/generic/README.md", "README for {{ .Project.Name }}")
	fsMock.AddTemplate("skeletons/postgresql/docker-compose.yml", `
services:
  db:
    image: postgres
    environment:
      POSTGRES_DB: {{ .Project.Name | ToLower }}_db
      POSTGRES_USER: {{ .Project.Name | ToLower }}_user
      POSTGRES_PASSWORD: {{ .Project.Name | ToLower }}_password
`)
	fsMock.AddTemplate("skeletons/generic/manifest.json", `{"files": [{"path": "README.md", "strategy": "overwrite"}]}`)
	fsMock.AddTemplate("skeletons/postgresql/manifest.json", `{"files": [{"path": "docker-compose.yml", "strategy": "merge-yaml"}]}`)

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
	err := service.Scaffold(tmpDir, fsMock.TempDir(), projRecipe)

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
