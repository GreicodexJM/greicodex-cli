package templates

import (
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	data := Data{
		ProjectName: "TestProject",
		Year:        2023,
	}

	content, err := Process("README.md.tmpl", data)
	if err != nil {
		t.Errorf("Process() returned an unexpected error: %v", err)
	}

	if !strings.Contains(string(content), "TestProject") {
		t.Errorf("Expected content to contain 'TestProject', but it did not")
	}

	if !strings.Contains(string(content), "2023") {
		t.Errorf("Expected content to contain '2023', but it did not")
	}
}

func TestProcess_MissingTemplate(t *testing.T) {
	data := Data{}
	_, err := Process("non-existent-template.tmpl", data)
	if err == nil {
		t.Error("Process() should have returned an error, but it did not")
	}
}
