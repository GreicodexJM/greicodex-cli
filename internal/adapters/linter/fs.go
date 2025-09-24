package linter

import (
	"fmt"
	"grei-cli/internal/ports/outbound"
	"os"
	"path/filepath"
)

// A map of linter names to their common configuration filenames.
var linterConfigFiles = map[string]string{
	"golangci-lint": ".golangci.yml",
	"ESLint":        ".eslintrc.js",
	"Ruff":          "pyproject.toml",
	"PHPStan":       "phpstan.neon",
}

type fsDetector struct{}

func NewFsDetector() outbound.LinterDetector {
	return &fsDetector{}
}

func (d *fsDetector) CheckConfig(path, linterName string) (bool, error) {
	configFile, ok := linterConfigFiles[linterName]
	if !ok {
		return false, fmt.Errorf("unknown linter: %s", linterName)
	}

	fullPath := filepath.Join(path, configFile)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
