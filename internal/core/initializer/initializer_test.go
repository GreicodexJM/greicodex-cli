package initializer

import (
	"errors"
	"grei-cli/internal/ports/outbound"
	"os"
	"testing"
)

type mockFSRepo struct {
	outbound.FSRepository
	createDirErr    error
	createFileErr   error
	cacheDir        string
	cacheDirErr     error
	readFileContent []byte
	readFileErr     error
}

func (m *mockFSRepo) CreateDir(path string) error {
	return m.createDirErr
}

func (m *mockFSRepo) CreateFile(path string, content []byte) error {
	return m.createFileErr
}

func (m *mockFSRepo) GetCacheDir(path string) (string, error) {
	return m.cacheDir, m.cacheDirErr
}

func (m *mockFSRepo) ReadFile(path string) ([]byte, error) {
	return m.readFileContent, m.readFileErr
}

func (m *mockFSRepo) ReadDir(path string) ([]os.DirEntry, error) {
	return nil, nil
}

type mockGitRepo struct {
	outbound.GitRepository
	initErr         error
	createBranchErr error
}

func (m *mockGitRepo) Init(path string) error {
	return m.initErr
}

func (m *mockGitRepo) CreateBranch(path, branchName string) error {
	return m.createBranchErr
}

func TestNewService(t *testing.T) {
	fsRepo := &mockFSRepo{}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)
	if service == nil {
		t.Error("NewService() should not return nil")
	}
}

func TestInitializeProject(t *testing.T) {
	fsRepo := &mockFSRepo{
		readFileContent: []byte(`{"minVersion": "0.1.0"}`),
	}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err != nil {
		t.Errorf("InitializeProject() returned an unexpected error: %v", err)
	}
}

func TestInitializeProject_GetCacheDirError(t *testing.T) {
	fsRepo := &mockFSRepo{cacheDirErr: errors.New("cache dir error")}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err == nil {
		t.Error("InitializeProject() should have returned an error, but it did not")
	}
}

func TestInitializeProject_DownloadError(t *testing.T) {
	fsRepo := &mockFSRepo{
		readFileContent: []byte(`{"minVersion": "0.1.0"}`),
	}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err != nil {
		t.Errorf("InitializeProject() returned an unexpected error: %v", err)
	}
}

func TestInitializeProject_CreateDirError(t *testing.T) {
	fsRepo := &mockFSRepo{createDirErr: errors.New("create dir error")}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err == nil {
		t.Error("InitializeProject() should have returned an error, but it did not")
	}
}

func TestInitializeProject_CreateFileError(t *testing.T) {
	fsRepo := &mockFSRepo{createFileErr: errors.New("create file error")}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err == nil {
		t.Error("InitializeProject() should have returned an error, but it did not")
	}
}

func TestInitializeProject_GitInitError(t *testing.T) {
	fsRepo := &mockFSRepo{}
	gitRepo := &mockGitRepo{initErr: errors.New("git init error")}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err == nil {
		t.Error("InitializeProject() should have returned an error, but it did not")
	}
}

func TestInitializeProject_CreateBranchError(t *testing.T) {
	fsRepo := &mockFSRepo{}
	gitRepo := &mockGitRepo{createBranchErr: errors.New("create branch error")}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", true)
	if err == nil {
		t.Error("InitializeProject() should have returned an error, but it did not")
	}
}

func TestInitializeProject_NoGitInit(t *testing.T) {
	fsRepo := &mockFSRepo{
		readFileContent: []byte(`{"minVersion": "0.1.0"}`),
	}
	gitRepo := &mockGitRepo{}
	service := NewService(fsRepo, gitRepo)

	err := service.InitializeProject("/tmp/test-project", "/tmp/cache", false)
	if err != nil {
		t.Errorf("InitializeProject() returned an unexpected error: %v", err)
	}
}
