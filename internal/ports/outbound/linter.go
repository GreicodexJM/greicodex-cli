package outbound

// LinterDetector defines the port for a service that can detect linter configurations.
type LinterDetector interface {
	// CheckConfig checks if the configuration file for a given linter exists.
	CheckConfig(path, linterName string) (bool, error)
}
