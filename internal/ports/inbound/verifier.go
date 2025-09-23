package inbound

type VerifyOptions struct {
	Path        string
	MinCoverage int
	JSONOutput  bool
}

// ProjectVerifier defines the port for the project verification service.
type ProjectVerifier interface {
	VerifyProject(options VerifyOptions) error
}
