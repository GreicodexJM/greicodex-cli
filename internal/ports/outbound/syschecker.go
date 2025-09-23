package outbound

// SystemChecker defines the port for checking system commands.
type SystemChecker interface {
	CommandExists(cmd string) bool
}
