package outbound

import "grei-cli/internal/core/plugin"

// PluginScanner defines the interface for discovering and loading plugins.
type PluginScanner interface {
	// Scan discovers all available plugins and returns their manifests.
	Scan() ([]*plugin.Manifest, error)
}
