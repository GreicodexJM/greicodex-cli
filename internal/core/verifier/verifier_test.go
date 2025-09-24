package verifier

import (
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"os"
	"path/filepath"
	"testing"
)

// Mock implementations of the outbound ports for testing.

type mockCoverageParser struct{}

func (m *mockCoverageParser) Parse(path string) (float64, error) {
	return 85.0, nil
}

type mockSecretScanner struct{}

func (m *mockSecretScanner) Scan(path string) ([]string, error) {
	return nil, nil
}

type mockLinterDetector struct{}

func (m *mockLinterDetector) CheckConfig(path, linterName string) (bool, error) {
	return true, nil
}

func TestVerifyProject_Success(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create dummy files that the verifier expects to find.
	dummyFiles := []string{
		"docker-compose.yml",
		"LICENSE",
		"CONTRIBUTING.md",
		"deploy/helm/Chart.yaml", // A file inside the dir is enough.
		"coverage.out",
	}
	for _, f := range dummyFiles {
		p := filepath.Join(tmpDir, f)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatalf("Failed to create parent dir for %s: %v", f, err)
		}
		if _, err := os.Create(p); err != nil {
			t.Fatalf("Failed to create dummy file %s: %v", f, err)
		}
	}

	coverageParsers := map[string]outbound.CoverageParser{
		"coverage.out": &mockCoverageParser{},
	}
	secretScanner := &mockSecretScanner{}
	linterDetector := &mockLinterDetector{}

	service := NewService(coverageParsers, secretScanner, linterDetector)

	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe: &recipe.Recipe{
			Project:     recipe.Project{Name: "TestProject"},
			Stack:       recipe.Stack{Linter: "golangci-lint"},
			Persistence: recipe.Persistence{Type: "postgresql"},
			Deployment:  recipe.Deployment{Type: "kubernetes"},
		},
	}

	// Act
	err = service.VerifyProject(options)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
