package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MockRepository struct {
	files   map[string][]byte
	tempDir string
}

func NewMockRepository() *MockRepository {
	tempDir, err := os.MkdirTemp("", "grei-test-*")
	if err != nil {
		panic(err)
	}
	return &MockRepository{
		files:   make(map[string][]byte),
		tempDir: tempDir,
	}
}

func (m *MockRepository) GetCacheDir(path string) (string, error) {
	return m.tempDir, nil
}

func (m *MockRepository) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (m *MockRepository) CreateFile(path string, content []byte) error {
	m.files[path] = content
	return os.WriteFile(path, content, 0644)
}

func (m *MockRepository) ReadFile(path string) ([]byte, error) {
	if content, exists := m.files[path]; exists {
		return content, nil
	}
	return os.ReadFile(path)
}

func (m *MockRepository) TempDir() string {
	return m.tempDir
}

func (m *MockRepository) Clean() {
	os.RemoveAll(m.tempDir)
}

func (m *MockRepository) AddFile(path string, content []byte) {
	fullPath := filepath.Join(m.tempDir, path)
	m.files[fullPath] = content
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		panic(err)
	}
}

func (m *MockRepository) AddTemplate(path string, content string) {
	fullPath := filepath.Join(m.tempDir, "templates", path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		panic(err)
	}
}

func (m *MockRepository) AddManifest(templateType, name, language, tooling, persistence string) {
	manifest := fmt.Sprintf(`
name: %s
description: A test template
type: %s
provides:
  language: %s
  tooling: %s
  persistence: %s
`, name, templateType, language, tooling, persistence)
	path := filepath.Join(name, "manifest.yml")
	if strings.Contains(name, "generic") {
		path = "manifest.yml"
	}
	m.AddTemplate(path, manifest)
}
