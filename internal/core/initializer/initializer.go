package initializer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"grei-cli/internal/templates"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Masterminds/semver"
)

const (
	templatesURL    = "https://github.com/GreicodexJM/greicodex-cli.git"
	templatesBranch = "master"
	templatesDir    = ".grei-cli/templates"
	cliVersion      = "0.1.0" // This should be replaced with a dynamic version
)

type Manifest struct {
	MinVersion string `json:"minVersion"`
}

type service struct {
	fsRepo     outbound.FSRepository
	gitRepo    outbound.GitRepository
	downloader outbound.Downloader
}

func NewService(fsRepo outbound.FSRepository, gitRepo outbound.GitRepository, downloader outbound.Downloader) inbound.InitializerService {
	return &service{
		fsRepo:     fsRepo,
		gitRepo:    gitRepo,
		downloader: downloader,
	}
}

func (s *service) InitializeProject(path string, gitInit bool) error {
	cacheDir, err := s.fsRepo.GetCacheDir(templatesDir)
	if err != nil {
		return err
	}

	if os.Getenv("GREI_E2E_TEST") == "" {
		if err := s.downloader.Download(context.Background(), templatesURL, templatesBranch, cacheDir); err != nil {
			fmt.Printf("warn: failed to download templates: %v. Using cached version if available.\n", err)
		}
	}

	if err := s.checkVersion(cacheDir); err != nil {
		return err
	}

	if err := s.fsRepo.CreateDir(path); err != nil {
		return err
	}

	projectName := filepath.Base(path)
	data := templates.Data{
		ProjectName: projectName,
		Year:        time.Now().Year(),
	}

	filesToCreate := map[string]string{
		"README.md":          "README.md.tmpl",
		".gitignore":         ".gitignore.tmpl",
		".editorconfig":      ".editorconfig.tmpl",
		"LICENSE":            "LICENSE.tmpl",
		"CONTRIBUTING.md":    "CONTRIBUTING.md.tmpl",
		"docker-compose.yml": "docker-compose.yml.tmpl",
	}

	for dest, tmpl := range filesToCreate {
		templatePath := filepath.Join(cacheDir, "templates", tmpl)
		content, err := s.fsRepo.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}

		t, err := template.New(tmpl).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var processedContent bytes.Buffer
		if err := t.Execute(&processedContent, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		err = s.fsRepo.CreateFile(filepath.Join(path, dest), processedContent.Bytes())
		if err != nil {
			return err
		}
	}

	dirsToCreate := []string{
		"docs",
		"deploy/helm",
		"docs/adr",
	}

	for _, dir := range dirsToCreate {
		err := s.fsRepo.CreateDir(filepath.Join(path, dir))
		if err != nil {
			return err
		}
	}

	if gitInit {
		if err := s.gitRepo.Init(path); err != nil {
			return err
		}
		if err := s.gitRepo.CreateBranch(path, "develop"); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) checkVersion(cacheDir string) error {
	manifestPath := filepath.Join(cacheDir, "templates", "manifest.json")
	manifestFile, err := s.fsRepo.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestFile, &manifest); err != nil {
		return fmt.Errorf("failed to unmarshal manifest file: %w", err)
	}

	minVersion, err := semver.NewVersion(manifest.MinVersion)
	if err != nil {
		return fmt.Errorf("failed to parse minVersion: %w", err)
	}

	currentVersion, err := semver.NewVersion(cliVersion)
	if err != nil {
		return fmt.Errorf("failed to parse cliVersion: %w", err)
	}

	if currentVersion.LessThan(minVersion) {
		return fmt.Errorf("cli version %s is less than the required minimum version %s", cliVersion, manifest.MinVersion)
	}

	return nil
}
