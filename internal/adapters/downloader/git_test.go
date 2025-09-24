package downloader

import (
	"context"
	"os"
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

	// Test pulling
	err = downloader.Download(context.Background(), "https://github.com/GreicodexJM/greicodex-cli.git", "master", cacheDir)
	if err != nil {
		t.Errorf("Download() failed to pull: %v", err)
	}
}
