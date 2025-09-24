package outbound

// FSRepository defines the port for file system operations.
type FSRepository interface {
	CreateDir(path string) error
	CreateFile(path string, content []byte) error
	GetCacheDir(path string) (string, error)
	ReadFile(path string) ([]byte, error)
}
