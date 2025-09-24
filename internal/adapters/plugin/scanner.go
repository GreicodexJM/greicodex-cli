package plugin

import (
	"encoding/json"
	"fmt"
	"grei-cli/internal/core/plugin"
	"grei-cli/internal/ports/outbound"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type scanner struct{}

// NewScanner creates a new plugin scanner.
func NewScanner() outbound.PluginScanner {
	return &scanner{}
}

// Scan discovers all available plugins (both built-in and external) and returns their manifests.
func (s *scanner) Scan() ([]*plugin.Manifest, error) {
	// 1. Load built-in plugins
	manifests := getBuiltInManifests()

	// 2. Scan for external plugins
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Non-fatal, we can still proceed with built-in plugins
		fmt.Fprintf(os.Stderr, "[Warning] Could not get user config directory: %v\n", err)
		return manifests, nil
	}

	pluginDir := filepath.Join(configDir, "grei", "plugins")
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// Directory doesn't exist, so no external plugins to scan.
		return manifests, nil
	}

	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, fmt.Errorf("could not read plugin directory '%s': %w", pluginDir, err)
	}

	for _, file := range files {
		if file.IsDir() || !isExecutable(file.Type()) {
			continue
		}

		pluginPath := filepath.Join(pluginDir, file.Name())
		if !strings.HasPrefix(file.Name(), "grei-") {
			fmt.Fprintf(os.Stderr, "[Warning] Found executable in plugin directory without 'grei-' prefix, skipping: %s\n", file.Name())
			continue
		}

		cmd := exec.Command(pluginPath, "discover")
		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Warning] Failed to execute 'discover' for plugin '%s': %v\n", file.Name(), err)
			continue
		}

		var manifest plugin.Manifest
		if err := json.Unmarshal(output, &manifest); err != nil {
			fmt.Fprintf(os.Stderr, "[Warning] Failed to parse manifest for plugin '%s': %v\n", file.Name(), err)
			continue
		}

		manifests = append(manifests, &manifest)
	}

	return manifests, nil
}

func isExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}

// getBuiltInManifests defines the plugins that are included with the CLI.
func getBuiltInManifests() []*plugin.Manifest {
	return []*plugin.Manifest{
		{
			Name:        "builtin-symfony-lamp",
			Description: "Pila tecnológica LAMP para aplicaciones web y APIs con Symfony.",
			Provides: []plugin.TechVertical{
				{
					ProjectType: "Aplicación Web",
					Stack: plugin.Stack{
						Language:  "PHP",
						Framework: "Symfony",
						WebServer: "Apache",
						Database:  "MySQL",
					},
				},
				{
					ProjectType: "API / Backend",
					Stack: plugin.Stack{
						Language:  "PHP",
						Framework: "API Platform (Symfony)",
						WebServer: "Apache",
						Database:  "MySQL",
					},
				},
			},
		},
		{
			Name:        "builtin-mern",
			Description: "Pila tecnológica MERN para aplicaciones web full-stack.",
			Provides: []plugin.TechVertical{
				{
					ProjectType: "Aplicación Web",
					Stack: plugin.Stack{
						Language:  "TypeScript",
						Framework: "React",
						WebServer: "Node.js",
						Database:  "MongoDB",
					},
				},
				{
					ProjectType: "API / Backend",
					Stack: plugin.Stack{
						Language:  "TypeScript",
						Framework: "Express",
						WebServer: "Node.js",
						Database:  "MongoDB",
					},
				},
			},
		},
	}
}
