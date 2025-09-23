package inbound

// ProjectInitializer defines the port for the project initialization service.
type ProjectInitializer interface {
	InitializeProject(path string, gitInit bool) error
}
