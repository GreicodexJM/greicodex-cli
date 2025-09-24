package linter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFsDetector(t *testing.T) {
	detector := NewFsDetector()
	if detector == nil {
		t.Error("NewFsDetector() should not return nil")
	}
}

func TestCheckConfig(t *testing.T) {
	detector := NewFsDetector()
	tmpDir, err := os.MkdirTemp("", "test-linter-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test for a linter config that exists
	configFilePath := filepath.Join(tmpDir, ".golangci.yml")
	if err := os.WriteFile(configFilePath, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	found, err := detector.CheckConfig(tmpDir, "golangci-lint")
	if err != nil {
		t.Errorf("CheckConfig() returned an unexpected error: %v", err)
	}
	if !found {
		t.Error("CheckConfig() should have returned true, but it did not")
	}

	// Test for a linter config that does not exist
	found, err = detector.CheckConfig(tmpDir, "ESLint")
	if err != nil {
		t.Errorf("CheckConfig() returned an unexpected error: %v", err)
	}
	if found {
		t.Error("CheckConfig() should have returned false, but it did not")
	}

	// Test for an unknown linter
	_, err = detector.CheckConfig(tmpDir, "unknown-linter")
	if err == nil {
		t.Error("CheckConfig() should have returned an error, but it did not")
	}

	// Test for a stat error
	readOnlyFilePath := filepath.Join(tmpDir, "read-only-file")
	if err := os.WriteFile(readOnlyFilePath, []byte{}, 0444); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_, err = detector.CheckConfig(readOnlyFilePath, "golangci-lint")
	if err == nil {
		t.Error("CheckConfig() should have returned an error, but it did not")
	}
}
