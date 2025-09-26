package inbound

// InitializerService defines the port for the project initialization service.
type InitializerService interface {
	InitializeProject(path, cacheDir string, gitInit bool) error
}
