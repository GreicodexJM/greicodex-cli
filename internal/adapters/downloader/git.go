package downloader

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitDownloader is an adapter for downloading files from a Git repository.
type GitDownloader struct{}

// NewGitDownloader creates a new GitDownloader.
func NewGitDownloader() *GitDownloader {
	return &GitDownloader{}
}

// Download clones or pulls a remote repository into a local cache directory.
func (d *GitDownloader) Download(ctx context.Context, url, branch, cacheDir string) error {
	_, err := os.Stat(cacheDir)
	if os.IsNotExist(err) {
		return d.clone(ctx, url, branch, cacheDir)
	}
	return d.pull(ctx, branch, cacheDir)
}

func (d *GitDownloader) clone(ctx context.Context, url, branch, cacheDir string) error {
	_, err := git.PlainCloneContext(ctx, cacheDir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

func (d *GitDownloader) pull(ctx context.Context, branch, cacheDir string) error {
	repo, err := git.PlainOpen(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.PullContext(ctx, &git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull repository: %w", err)
	}

	return nil
}
