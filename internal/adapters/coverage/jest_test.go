package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewJestParser(t *testing.T) {
	parser := NewJestParser()
	if parser == nil {
		t.Error("NewJestParser() should not return nil")
	}
}

func TestParse(t *testing.T) {
	parser := NewJestParser()
	tmpDir, err := os.MkdirTemp("", "test-coverage-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test with a valid coverage file
	coverageFilePath := filepath.Join(tmpDir, "coverage-summary.json")
	coverageFileContent := `{
		"total": {
			"lines": {
				"pct": 85.5
			}
		}
	}`
	if err := os.WriteFile(coverageFilePath, []byte(coverageFileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	pct, err := parser.Parse(coverageFilePath)
	if err != nil {
		t.Errorf("Parse() returned an unexpected error: %v", err)
	}
	if pct != 85.5 {
		t.Errorf("Expected coverage to be 85.5, but got %f", pct)
	}

	// Test with a missing file
	_, err = parser.Parse(filepath.Join(tmpDir, "non-existent-file.json"))
	if err == nil {
		t.Error("Parse() should have returned an error, but it did not")
	}

	// Test with an invalid JSON file
	invalidJSONFilePath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidJSONFilePath, []byte("{"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_, err = parser.Parse(invalidJSONFilePath)
	if err == nil {
		t.Error("Parse() should have returned an error, but it did not")
	}
}
