package filesystem

import (
	"grei-cli/internal/ports/outbound"
	"os"
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
