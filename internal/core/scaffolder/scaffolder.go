package scaffolder

import (
	"bytes"
	"embed"
	"fmt"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/ports/inbound"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:embed all:templates
var templateFiles embed.FS

func GetTemplates() ([]fs.DirEntry, error) {
	return templateFiles.ReadDir("templates")
}

func GetTemplateFile(path string) ([]byte, error) {
	return templateFiles.ReadFile(path)
}

type service struct{}

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

func NewService() inbound.ScaffolderService {
	return &service{}
}

func (s *service) Scaffold(path string, recipe *recipe.Recipe) error {
	fmt.Printf("\n[i] Scaffolding templates for a '%s' project...\n", recipe.Project.Type)

	// 1. Copy generic templates
	if err := s.copyTemplates(filepath.Join("templates", "generic"), path, recipe); err != nil {
		return err
	}

	// 2. Copy stack-specific templates
	templateDirs, err := fs.ReadDir(templateFiles, "templates")
	if err != nil {
		return err
	}

	for _, dir := range templateDirs {
		if !dir.IsDir() || dir.Name() == "generic" {
			continue
		}

		manifestPath := filepath.Join("templates", dir.Name(), "manifest.yml")
		manifestFile, err := templateFiles.ReadFile(manifestPath)
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
			if err := s.copyTemplates(filepath.Join("templates", dir.Name()), path, recipe); err != nil {
				return err
			}
		}

		if manifest.Provides.Persistence == recipe.Persistence.Type {
			if err := s.copyTemplates(filepath.Join("templates", dir.Name()), path, recipe); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) copyTemplates(sourceDir, targetDir string, recipe *recipe.Recipe) error {
	return fs.WalkDir(templateFiles, sourceDir, func(templatePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// The path in the target directory.
		relativePath := strings.TrimPrefix(templatePath, sourceDir)
		if strings.HasSuffix(relativePath, ".tmpl") {
			relativePath = strings.TrimSuffix(relativePath, ".tmpl")
		}
		targetPath := filepath.Join(targetDir, relativePath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read the template file.
		rawContent, err := templateFiles.ReadFile(templatePath)
		if err != nil {
			return err
		}

		// Execute the template to replace variables like {{ .Project.Name }}
		funcMap := template.FuncMap{
			"ToLower": strings.ToLower,
		}
		tmpl, err := template.New(d.Name()).Funcs(funcMap).Parse(string(rawContent))
		if err != nil {
			return fmt.Errorf("could not parse template %s: %w", templatePath, err)
		}

		var processedContent bytes.Buffer
		if err := tmpl.Execute(&processedContent, recipe); err != nil {
			return fmt.Errorf("could not execute template %s: %w", templatePath, err)
		}

		// Write the file to the target directory.
		return os.WriteFile(targetPath, processedContent.Bytes(), 0644)
	})
}
