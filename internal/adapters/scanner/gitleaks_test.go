package scanner

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

func TestNewGitleaksScanner(t *testing.T) {
	sysChecker := &mockSysChecker{}
	scanner := NewGitleaksScanner(sysChecker)
	if scanner == nil {
		t.Error("NewGitleaksScanner() should not return nil")
	}
}

func TestScan_GitleaksNotFound(t *testing.T) {
	sysChecker := &mockSysChecker{commandExists: false}
	scanner := NewGitleaksScanner(sysChecker)

	_, err := scanner.Scan("/tmp")
	if err != ErrGitleaksNotFound {
		t.Errorf("Scan() should have returned ErrGitleaksNotFound, but it did not")
	}
}

func TestScan_GitleaksError(t *testing.T) {
	sysChecker := &mockSysChecker{commandExists: true}
	scanner := NewGitleaksScanner(sysChecker)

	_, err := scanner.Scan("/non-existent-path")
	if err == nil {
		t.Error("Scan() should have returned an error, but it did not")
	}
}
