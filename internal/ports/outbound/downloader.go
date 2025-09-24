package outbound

import "context"

// Downloader is the port for downloading files from a remote repository.
type Downloader interface {
	// Download clones or pulls a remote repository into a local cache directory.
	Download(ctx context.Context, url, branch, cacheDir string) error
}
