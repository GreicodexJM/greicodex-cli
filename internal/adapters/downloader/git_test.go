package downloader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGitDownloader_Download(t *testing.T) {
	// Create a temporary directory for the cache
	cacheDir, err := os.MkdirTemp("", "grei-cli-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(cacheDir)

	downloader := NewGitDownloader()

	// Test cloning
	err = downloader.Download(context.Background(), "https://github.com/GreicodexJM/greicodex-cli.git", "master", cacheDir)
	if err != nil {
		t.Errorf("Download() failed to clone: %v", err)
	}

	// Verify that only the templates folder was downloaded
	_, err = os.Stat(filepath.Join(cacheDir, "templates"))
	if os.IsNotExist(err) {
		t.Error("templates directory was not downloaded")
	}

	_, err = os.Stat(filepath.Join(cacheDir, "README.md"))
	if !os.IsNotExist(err) {
		t.Error("README.md was downloaded, but it should not have been")
	}

	// Test pulling
	err = downloader.Download(context.Background(), "https://github.com/GreicodexJM/greicodex-cli.git", "master", cacheDir)
	if err != nil {
		t.Errorf("Download() failed to pull: %v", err)
	}
}
