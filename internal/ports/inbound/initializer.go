package inbound

// InitializerService defines the port for the project initialization service.
type InitializerService interface {
	InitializeProject(path string, gitInit bool) error
}
