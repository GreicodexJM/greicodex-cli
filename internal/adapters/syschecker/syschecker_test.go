package syschecker

import (
	"testing"
)

func TestNew(t *testing.T) {
	checker := New()
	if checker == nil {
		t.Error("New() should not return nil")
	}
}

func TestCommandExists(t *testing.T) {
	checker := New()

	// Test for a command that should exist
	if !checker.CommandExists("go") {
		t.Error("CommandExists('go') should have returned true, but it did not")
	}

	// Test for a command that should not exist
	if checker.CommandExists("non-existent-command") {
		t.Error("CommandExists('non-existent-command') should have returned false, but it did not")
	}
}
