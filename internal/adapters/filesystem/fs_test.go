package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	if repo == nil {
		t.Error("NewRepository() should not return nil")
	}
}

func TestCreateDir(t *testing.T) {
	repo := NewRepository()
	tmpDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dirPath := filepath.Join(tmpDir, "new-dir")
	err = repo.CreateDir(dirPath)
	if err != nil {
		t.Errorf("CreateDir() returned an unexpected error: %v", err)
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Expected directory to be created, but it was not: %s", dirPath)
	}
}

func TestCreateFile(t *testing.T) {
	repo := NewRepository()
	tmpDir, err := os.MkdirTemp("", "test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "new-file.txt")
	content := []byte("hello world")
	err = repo.CreateFile(filePath, content)
	if err != nil {
		t.Errorf("CreateFile() returned an unexpected error: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file to be created, but it was not: %s", filePath)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(fileContent) != string(content) {
		t.Errorf("Expected file content to be '%s', but got '%s'", string(content), string(fileContent))
	}
}
