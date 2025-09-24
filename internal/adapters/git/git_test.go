package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	if repo == nil {
		t.Error("NewRepository() should not return nil")
	}
}

func TestInit(t *testing.T) {
	repo := NewRepository()
	tmpDir, err := os.MkdirTemp("", "test-git-init-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = repo.Init(tmpDir)
	if err != nil {
		t.Errorf("Init() returned an unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".git")); os.IsNotExist(err) {
		t.Errorf("Expected .git directory to be created, but it was not")
	}
}

func TestSetConfig(t *testing.T) {
	repo := NewRepository()
	tmpDir, err := os.MkdirTemp("", "test-git-config-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = repo.Init(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	err = repo.SetConfig(tmpDir, "user.name", "Test User")
	if err != nil {
		t.Errorf("SetConfig() returned an unexpected error: %v", err)
	}

	cmd := exec.Command("git", "config", "user.name")
	cmd.Dir = tmpDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get git config: %v", err)
	}

	if string(out) != "Test User\n" {
		t.Errorf("Expected git config user.name to be 'Test User', but got '%s'", string(out))
	}
}

func TestCreateBranch(t *testing.T) {
	repo := NewRepository()
	tmpDir, err := os.MkdirTemp("", "test-git-branch-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = repo.Init(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Need to commit something before creating a new branch
	err = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}

	err = repo.CreateBranch(tmpDir, "new-branch")
	if err != nil {
		t.Errorf("CreateBranch() returned an unexpected error: %v", err)
	}

	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = tmpDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	if string(out) != "new-branch\n" {
		t.Errorf("Expected current branch to be 'new-branch', but got '%s'", string(out))
	}
}
