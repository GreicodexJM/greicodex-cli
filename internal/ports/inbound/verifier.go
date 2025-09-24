package inbound

import "grei-cli/internal/core/recipe"

type VerifyOptions struct {
	Path        string
	MinCoverage int
	JSONOutput  bool
	Recipe      *recipe.Recipe
}

// ProjectVerifier defines the port for the project verification service.
type ProjectVerifier interface {
	VerifyProject(options VerifyOptions) error
}
