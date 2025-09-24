package inbound

import "grei-cli/internal/core/recipe"

type VerifyOptions struct {
	Path        string
	MinCoverage int
	JSONOutput  bool
	Recipe      *recipe.Recipe
}

// VerifierService defines the port for the project verification service.
type VerifierService interface {
	VerifyProject(options VerifyOptions) error
}
