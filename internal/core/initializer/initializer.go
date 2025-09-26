package initializer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"grei-cli/internal/templates"
	"path/filepath"
	"strings"
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
	fsRepo  outbound.FSRepository
	gitRepo outbound.GitRepository
}

func NewService(fsRepo outbound.FSRepository, gitRepo outbound.GitRepository) inbound.InitializerService {
	return &service{
		fsRepo:  fsRepo,
		gitRepo: gitRepo,
	}
}

func (s *service) InitializeProject(path, cacheDir string, gitInit bool, recipe *recipe.Recipe) error {
	templatesCacheDir := filepath.Join(cacheDir, "templates")
	if err := s.checkVersion(templatesCacheDir); err != nil {
		return err
	}

	if err := s.fsRepo.CreateDir(path); err != nil {
		return err
	}

	data := templates.Data{
		Recipe: *recipe,
		Year:   time.Now().Year(),
	}

	genericSkeletonPath := filepath.Join(cacheDir, "templates", "skeletons", "generic")
	files, err := s.fsRepo.ReadDir(genericSkeletonPath)
	if err != nil {
		return fmt.Errorf("failed to read generic skeleton directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		templatePath := filepath.Join(genericSkeletonPath, file.Name())
		content, err := s.fsRepo.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}

		t, err := template.New(file.Name()).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var processedContent bytes.Buffer
		if err := t.Execute(&processedContent, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		destFileName := strings.TrimSuffix(file.Name(), ".tmpl")
		err = s.fsRepo.CreateFile(filepath.Join(path, destFileName), processedContent.Bytes())
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
	manifestPath := filepath.Join(cacheDir, "manifest.json")
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
