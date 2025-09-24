package doctor

import (
	"grei-cli/internal/ports/outbound"
	"testing"
)

type mockSysChecker struct {
	outbound.SystemChecker
	commandExists bool
}

func (m *mockSysChecker) CommandExists(command string) bool {
	return m.commandExists
}

func TestNewService(t *testing.T) {
	sysChecker := &mockSysChecker{}
	service := NewService(sysChecker)
	if service == nil {
		t.Error("NewService() should not return nil")
	}
}

func TestCheckEnvironment(t *testing.T) {
	sysChecker := &mockSysChecker{commandExists: true}
	service := NewService(sysChecker)

	results := service.CheckEnvironment()
	if len(results) != len(requiredTools) {
		t.Errorf("Expected %d results, but got %d", len(requiredTools), len(results))
	}

	for _, result := range results {
		if !result.Found {
			t.Errorf("Expected command '%s' to be found, but it was not", result.Command)
		}
	}
}

func TestCheckEnvironment_CommandNotFound(t *testing.T) {
	sysChecker := &mockSysChecker{commandExists: false}
	service := NewService(sysChecker)

	results := service.CheckEnvironment()
	if len(results) != len(requiredTools) {
		t.Errorf("Expected %d results, but got %d", len(requiredTools), len(results))
	}

	for _, result := range results {
		if result.Found {
			t.Errorf("Expected command '%s' to not be found, but it was", result.Command)
		}
	}
}
