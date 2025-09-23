package git

import (
	"grei-cli/internal/ports/outbound"
	"os/exec"
)

type repository struct{}

func NewRepository() outbound.GitRepository {
	return &repository{}
}

func (r *repository) SetConfig(path, key, value string) error {
	cmd := exec.Command("git", "config", key, value)
	cmd.Dir = path
	return cmd.Run()
}

func (r *repository) Init(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	return cmd.Run()
}

func (r *repository) CreateBranch(path, branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = path
	return cmd.Run()
}
