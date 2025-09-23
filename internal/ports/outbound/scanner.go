package outbound

// SecretScanner defines the port for scanning for secrets.
type SecretScanner interface {
	Scan(path string) ([]string, error)
}
