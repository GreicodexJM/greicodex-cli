package hooks

import (
	"errors"
	"grei-cli/internal/ports/outbound"
	"testing"
)

type mockGitRepo struct {
	outbound.GitRepository
	setConfigErr error
}

func (m *mockGitRepo) SetConfig(path, key, value string) error {
	return m.setConfigErr
}

func TestNewService(t *testing.T) {
	gitRepo := &mockGitRepo{}
	service := NewService(gitRepo)
	if service == nil {
		t.Error("NewService() should not return nil")
	}
}

func TestInstallHooks(t *testing.T) {
	gitRepo := &mockGitRepo{}
	service := NewService(gitRepo)

	err := service.InstallHooks("/tmp/test-project")
	if err != nil {
		t.Errorf("InstallHooks() returned an unexpected error: %v", err)
	}
}

func TestInstallHooks_Error(t *testing.T) {
	gitRepo := &mockGitRepo{setConfigErr: errors.New("set config error")}
	service := NewService(gitRepo)

	err := service.InstallHooks("/tmp/test-project")
	if err == nil {
		t.Error("InstallHooks() should have returned an error, but it did not")
	}
}
