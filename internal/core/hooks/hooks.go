package hooks

import (
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
)

type service struct {
	gitRepo outbound.GitRepository
}

func NewService(gitRepo outbound.GitRepository) inbound.HookInstaller {
	return &service{
		gitRepo: gitRepo,
	}
}

func (s *service) InstallHooks(path string) error {
	return s.gitRepo.SetConfig(path, "core.hooksPath", ".githooks")
}
