package inbound

// HookInstaller defines the port for the Git hooks installation service.
type HookInstaller interface {
	InstallHooks(path string) error
}
