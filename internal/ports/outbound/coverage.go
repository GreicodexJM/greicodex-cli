package outbound

// CoverageParser defines the port for parsing test coverage reports.
type CoverageParser interface {
	Parse(path string) (float64, error)
}
