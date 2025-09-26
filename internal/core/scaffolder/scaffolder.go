package scaffolder

import (
	"bytes"
	"fmt"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	templatesDir = ".grei-cli/templates"
)

type service struct {
	fsRepo outbound.FSRepository
}

type Manifest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Provides    struct {
		Language             string `yaml:"language"`
		Tooling              string `yaml:"tooling"`
		Runtime              string `yaml:"runtime"`
		Persistence          string `yaml:"persistence"`
		DependencyManagement string `yaml:"dependencyManagement"`
		BuildReleaseRun      string `yaml:"buildReleaseRun"`
	} `yaml:"provides"`
}

func NewService(fsRepo outbound.FSRepository) inbound.ScaffolderService {
	return &service{
		fsRepo: fsRepo,
	}
}

func (s *service) GetTemplates() ([]fs.DirEntry, error) {
	cacheDir, err := s.fsRepo.GetCacheDir(templatesDir)
	if err != nil {
		return nil, err
	}
	return os.ReadDir(cacheDir)
}

func (s *service) GetTemplateFile(path string) ([]byte, error) {
	cacheDir, err := s.fsRepo.GetCacheDir(templatesDir)
	if err != nil {
		return nil, err
	}
	return s.fsRepo.ReadFile(filepath.Join(cacheDir, path))
}

func (s *service) Scaffold(path string, recipe *recipe.Recipe) error {
	fmt.Printf("\n[i] Scaffolding templates for a '%s' project...\n", recipe.Project.Type)

	cacheDir, err := s.fsRepo.GetCacheDir(templatesDir)
	if err != nil {
		return err
	}

	// 1. Copy generic templates
	if err := s.copyTemplates(filepath.Join(cacheDir, "generic"), path, recipe); err != nil {
		return err
	}

	// 2. Copy stack-specific templates
	templateDirs, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	for _, dir := range templateDirs {
		if !dir.IsDir() || dir.Name() == "generic" {
			continue
		}

		manifestPath := filepath.Join(cacheDir, dir.Name(), "manifest.yml")
		manifestFile, err := s.fsRepo.ReadFile(manifestPath)
		if err != nil {
			fmt.Printf("warn: could not read manifest for template %s: %v\n", dir.Name(), err)
			continue
		}

		var manifest Manifest
		if err := yaml.Unmarshal(manifestFile, &manifest); err != nil {
			fmt.Printf("warn: could not unmarshal manifest for template %s: %v\n", dir.Name(), err)
			continue
		}

		if manifest.Provides.Language == recipe.Stack.Language && manifest.Provides.Tooling == recipe.Stack.Tooling {
			if err := s.copyTemplates(filepath.Join(cacheDir, dir.Name()), path, recipe); err != nil {
				return err
			}
		}

		if manifest.Provides.Persistence == recipe.Persistence.Type {
			if err := s.copyTemplates(filepath.Join(cacheDir, dir.Name()), path, recipe); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) copyTemplates(sourceDir, targetDir string, recipe *recipe.Recipe) error {
	return filepath.Walk(sourceDir, func(templatePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// The path in the target directory.
		relativePath := strings.TrimPrefix(templatePath, sourceDir)
		if strings.HasSuffix(relativePath, ".tmpl") {
			relativePath = strings.TrimSuffix(relativePath, ".tmpl")
		}
		targetPath := filepath.Join(targetDir, relativePath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read the template file.
		rawContent, err := s.fsRepo.ReadFile(templatePath)
		if err != nil {
			return err
		}

		// Execute the template to replace variables like {{ .Project.Name }}
		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
		}
		tmpl, err := template.New(info.Name()).Funcs(funcMap).Parse(string(rawContent))
		if err != nil {
			return fmt.Errorf("could not parse template %s: %w", templatePath, err)
		}

		var processedContent bytes.Buffer
		if err := tmpl.Execute(&processedContent, recipe); err != nil {
			return fmt.Errorf("could not execute template %s: %w", templatePath, err)
		}

		// Write the file to the target directory.
		return s.fsRepo.CreateFile(targetPath, processedContent.Bytes())
	})
}
