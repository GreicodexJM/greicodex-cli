package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cliPath string
)

func TestMain(m *testing.M) {
	var err error
	cliPath, err = buildCLI()
	if err != nil {
		os.Exit(1)
	}
	exitCode := m.Run()
	cleanup(cliPath)
	os.Exit(exitCode)
}

func buildCLI() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("go", "build", "-o", "grei.exe", ".")
	} else {
		cmd = exec.Command("go", "build", "-o", "grei", ".")
	}
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return filepath.Abs(cmd.Args[3])
}

func cleanup(path string) {
	os.Remove(path)
}

func TestE2EInitCommand(t *testing.T) {
	t.Run("Successfully initialize a new project non-interactively", func(t *testing.T) {
		// Setup
		tempDir, err := os.MkdirTemp("", "grei-e2e-test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		recipePath := createTestRecipe(t, tempDir)

		// Execute
		cmd := exec.Command(cliPath, "init", "my-new-project", "--no-interactive", "--recipe-file", recipePath)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		// Assert
		assert.NoError(t, err, string(output))
		projectPath := filepath.Join(tempDir, "my-new-project")
		assert.FileExists(t, filepath.Join(projectPath, "grei.yml"))
		assert.FileExists(t, filepath.Join(projectPath, "README.md"))
		assert.FileExists(t, filepath.Join(projectPath, ".gitignore"))
		assert.DirExists(t, filepath.Join(projectPath, "docs"))
		assert.DirExists(t, filepath.Join(projectPath, ".git"))
	})

	t.Run("Fail if project already exists", func(t *testing.T) {
		// Setup
		tempDir, err := os.MkdirTemp("", "grei-e2e-test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a dummy grei.yml to trigger the error
		err = os.WriteFile(filepath.Join(tempDir, "grei.yml"), []byte("test"), 0644)
		assert.NoError(t, err)

		// Execute
		cmd := exec.Command(cliPath, "init", ".")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, string(output), "este directorio ya contiene un proyecto 'grei' (grei.yml encontrado)")
	})
}

func TestE2EVerifyCommand(t *testing.T) {
	t.Run("Successfully verify a compliant project", func(t *testing.T) {
		// Setup
		tempDir, err := os.MkdirTemp("", "grei-e2e-test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a project to verify
		cmd := exec.Command(cliPath, "init", "my-project", "--no-interactive", "--recipe-file", createTestRecipe(t, tempDir))
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(output))

		// Create a dummy coverage file
		projectPath := filepath.Join(tempDir, "my-project")
		err = os.WriteFile(filepath.Join(projectPath, "coverage.out"), []byte(`mode: set
github.com/user/project/file1.go:1.1,2.2 1 1
github.com/user/project/file2.go:1.1,2.2 1 1
`), 0644)
		assert.NoError(t, err)

		// Execute
		cmd = exec.Command(cliPath, "verify")
		cmd.Dir = filepath.Join(tempDir, "my-project")
		output, err = cmd.CombinedOutput()

		// Assert
		assert.NoError(t, err, string(output))
		assert.Contains(t, string(output), "¡Proyecto verificado exitosamente!")
	})
}

func TestE2EInstallHooksCommand(t *testing.T) {
	t.Run("Successfully install hooks", func(t *testing.T) {
		// Setup
		tempDir, err := os.MkdirTemp("", "grei-e2e-test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Initialize a git repository
		cmd := exec.Command("git", "init")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(output))

		// Execute
		cmd = exec.Command(cliPath, "install-hooks")
		cmd.Dir = tempDir
		output, err = cmd.CombinedOutput()

		// Assert
		assert.NoError(t, err, string(output))

		// Verify that the git config was updated
		configCmd := exec.Command("git", "config", "--get", "core.hooksPath")
		configCmd.Dir = tempDir
		configOutput, err := configCmd.CombinedOutput()
		assert.NoError(t, err, string(configOutput))
		assert.Equal(t, ".githooks\n", string(configOutput))
	})
}

func TestE2EDoctorCommand(t *testing.T) {
	t.Run("Successfully run doctor", func(t *testing.T) {
		// Execute
		cmd := exec.Command(cliPath, "doctor")
		output, err := cmd.CombinedOutput()

		// Assert
		assert.NoError(t, err, string(output))
		assert.Contains(t, string(output), "¡Tu entorno está listo!")
	})
}

func createTestRecipe(t *testing.T, dir string) string {
	recipe := `
project:
  name: my-project
  customer: Greicodex
  type: cli
`
	recipePath := filepath.Join(dir, "recipe.yml")
	err := os.WriteFile(recipePath, []byte(recipe), 0644)
	assert.NoError(t, err)
	return recipePath
}
