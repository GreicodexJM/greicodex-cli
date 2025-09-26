package inbound

import (
	"grei-cli/internal/core/recipe"
	"io/fs"
)

// ScaffolderService defines the port for the project scaffolding service.
type ScaffolderService interface {
	Scaffold(path, cacheDir string, recipe *recipe.Recipe) error
	GetTemplates() ([]fs.DirEntry, error)
	GetTemplateFile(path string) ([]byte, error)
}
