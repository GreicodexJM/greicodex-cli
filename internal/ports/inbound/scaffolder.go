package inbound

import "grei-cli/internal/core/recipe"

// Scaffolder defines the port for the project scaffolding service.
type Scaffolder interface {
	Scaffold(path string, recipe *recipe.Recipe) error
}
