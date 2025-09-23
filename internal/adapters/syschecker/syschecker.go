package syschecker

import (
	"grei-cli/internal/ports/outbound"
	"os/exec"
)

type checker struct{}

func New() outbound.SystemChecker {
	return &checker{}
}

func (c *checker) CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
