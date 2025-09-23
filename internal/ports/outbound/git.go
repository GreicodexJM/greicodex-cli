package outbound

// GitRepository defines the port for Git operations.
type GitRepository interface {
	SetConfig(path, key, value string) error
	Init(path string) error
	CreateBranch(path, branchName string) error
}
