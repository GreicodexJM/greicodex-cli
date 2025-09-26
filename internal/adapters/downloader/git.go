package downloader

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GitDownloader is an adapter for downloading files from a Git repository.
type GitDownloader struct{}

// NewGitDownloader creates a new GitDownloader.
func NewGitDownloader() *GitDownloader {
	return &GitDownloader{}
}

// Download clones or pulls a remote repository into a local cache directory.
func (d *GitDownloader) Download(ctx context.Context, url, branch, cacheDir string) error {
	gitDir := filepath.Join(cacheDir, ".git")
	_, err := os.Stat(gitDir)
	if os.IsNotExist(err) {
		return d.sparseClone(ctx, url, branch, cacheDir)
	}
	return d.pull(ctx, branch, cacheDir)
}

func (d *GitDownloader) sparseClone(ctx context.Context, url, branch, cacheDir string) error {
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Initialize an empty repository
	if err := d.runGitCommand(ctx, cacheDir, "init"); err != nil {
		return err
	}

	// Add the remote
	if err := d.runGitCommand(ctx, cacheDir, "remote", "add", "origin", url); err != nil {
		return err
	}

	// Enable sparse checkout
	if err := d.runGitCommand(ctx, cacheDir, "config", "core.sparseCheckout", "true"); err != nil {
		return err
	}

	// Define the sparse checkout directory
	sparseCheckoutFile := filepath.Join(cacheDir, ".git", "info", "sparse-checkout")
	if err := os.WriteFile(sparseCheckoutFile, []byte("templates"), 0644); err != nil {
		return fmt.Errorf("failed to write sparse-checkout file: %w", err)
	}

	// Pull the files
	if err := d.runGitCommand(ctx, cacheDir, "pull", "origin", branch); err != nil {
		return err
	}

	return nil
}

func (d *GitDownloader) pull(ctx context.Context, branch, cacheDir string) error {
	return d.runGitCommand(ctx, cacheDir, "pull", "origin", branch)
}

func (d *GitDownloader) runGitCommand(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git command %v: %w", args, err)
	}
	return nil
}
