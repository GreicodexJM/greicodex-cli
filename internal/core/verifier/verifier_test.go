package verifier

import (
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/ports/inbound"
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

type mockLinterDetector struct {
	checkConfigFunc func(path, linterName string) (bool, error)
}

func (m *mockLinterDetector) CheckConfig(path, linterName string) (bool, error) {
	if m.checkConfigFunc != nil {
		return m.checkConfigFunc(path, linterName)
	}
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

	coverageParser := &mockCoverageParser{}
	secretScanner := &mockSecretScanner{}
	linterDetector := &mockLinterDetector{}

	service := NewService(coverageParser, secretScanner, linterDetector)

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

func TestVerifyProject_CoverageTooLow(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	coverageParser := &mockCoverageParser{} // Returns 85.0
	service := NewService(coverageParser, &mockSecretScanner{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 90, // Higher than mock parser's return value
		Recipe:      &recipe.Recipe{},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for low coverage, but got none")
	}
}

func TestVerifyProject_MissingLinterConfig(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	linterDetector := &mockLinterDetector{checkConfigFunc: func(path, linterName string) (bool, error) {
		return false, nil // Simulate not found
	}}
	service := NewService(&mockCoverageParser{}, &mockSecretScanner{}, linterDetector)
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe:      &recipe.Recipe{Stack: recipe.Stack{Linter: "golangci-lint"}},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for missing linter config, but got none")
	}
}

func TestVerifyProject_MissingPersistence(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	service := NewService(&mockCoverageParser{}, &mockSecretScanner{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe: &recipe.Recipe{
			Persistence: recipe.Persistence{Type: "postgresql"},
		},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for missing persistence, but got none")
	}
}

func TestVerifyProject_MissingDeployment(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	service := NewService(&mockCoverageParser{}, &mockSecretScanner{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe: &recipe.Recipe{
			Deployment: recipe.Deployment{Type: "kubernetes"},
		},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for missing deployment, but got none")
	}
}

type mockSecretScannerSecretsFound struct{}

func (m *mockSecretScannerSecretsFound) Scan(path string) ([]string, error) {
	return []string{"secret found"}, nil
}

func TestVerifyProject_SecretsFound(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	service := NewService(&mockCoverageParser{}, &mockSecretScannerSecretsFound{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe:      &recipe.Recipe{},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for secrets found, but got none")
	}
}

func TestVerifyProject_MissingRequiredFiles(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)
	os.Create(filepath.Join(tmpDir, "coverage.out"))

	service := NewService(&mockCoverageParser{}, &mockSecretScanner{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe:      &recipe.Recipe{},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for missing required files, but got none")
	}
}

func TestVerifyProject_MissingCoverageFile(t *testing.T) {
	// Arrange
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)

	service := NewService(&mockCoverageParser{}, &mockSecretScanner{}, &mockLinterDetector{})
	options := inbound.VerifyOptions{
		Path:        tmpDir,
		MinCoverage: 80,
		Recipe:      &recipe.Recipe{},
	}

	// Act
	err := service.VerifyProject(options)

	// Assert
	if err == nil {
		t.Error("Expected an error for missing coverage file, but got none")
	}
}
