package inbound

// HooksService defines the port for the Git hooks installation service.
type HooksService interface {
	InstallHooks(path string) error
}
