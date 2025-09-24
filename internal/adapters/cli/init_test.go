package cli

import (
	"grei-cli/internal/core/stack"
	"strings"
	"testing"
)

func TestCategorizeStacks(t *testing.T) {
	// Reset the registry to a known state for this test.
	originalRegistry := stack.Registry
	stack.Registry = []*stack.Stack{
		{Name: "go-cli", Type: "code"},
		{Name: "postgresql", Type: "persistence"},
		{Name: "kubernetes", Type: "deployment"},
	}
	defer func() { stack.Registry = originalRegistry }()

	code, persistence, deployment := categorizeStacks()

	if len(code) != 2 || code[0] != "Custom" || code[1] != "go-cli" {
		t.Errorf("Expected code stacks to be [Custom go-cli], but got %v", code)
	}
	if len(persistence) != 2 || persistence[0] != "Ninguna" || persistence[1] != "postgresql" {
		t.Errorf("Expected persistence stacks to be [Ninguna postgresql], but got %v", persistence)
	}
	if len(deployment) != 2 || deployment[0] != "Ninguno" || deployment[1] != "kubernetes" {
		t.Errorf("Expected deployment stacks to be [Ninguno kubernetes], but got %v", deployment)
	}
}

func TestGenerateProjectName(t *testing.T) {
	name1 := generateProjectName()
	name2 := generateProjectName()

	if name1 == "" {
		t.Error("generateProjectName() should not return an empty string")
	}
	if name1 == name2 {
		t.Logf("Generated names were the same ('%s'), which is possible but unlikely. Running test again.", name1)
		name2 = generateProjectName()
		if name1 == name2 {
			t.Errorf("generateProjectName() returned the same name twice in a row ('%s'), which is highly unlikely.", name1)
		}
	}

	// Check that the name is in CamelCase
	if !strings.ContainsAny(name1, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("Expected project name to contain uppercase letters, but it did not: %s", name1)
	}
}
