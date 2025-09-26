package filesystem

import (
	"grei-cli/internal/ports/outbound"
	"os"
	"path/filepath"
)

type repository struct{}

func NewRepository() outbound.FSRepository {
	return &repository{}
}

func (r *repository) CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (r *repository) CreateFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}

func (r *repository) GetCacheDir(path string) (string, error) {
	if cacheDir := os.Getenv("GREI_CACHE_DIR"); cacheDir != "" {
		return filepath.Join(cacheDir, path), nil
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, path), nil
}

func (r *repository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
