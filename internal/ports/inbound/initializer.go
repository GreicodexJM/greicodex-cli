package inbound

import "grei-cli/internal/core/recipe"

// InitializerService defines the port for the project initialization service.
type InitializerService interface {
	InitializeProject(path, cacheDir string, gitInit bool, recipe *recipe.Recipe) error
}
