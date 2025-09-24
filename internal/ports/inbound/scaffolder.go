package inbound

import "grei-cli/internal/core/recipe"

// ScaffolderService defines the port for the project scaffolding service.
type ScaffolderService interface {
	Scaffold(path string, recipe *recipe.Recipe) error
}
