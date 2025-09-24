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
	"text/template"
)

//go:embed templates
var templateFiles embed.FS

type service struct{}

func NewService() inbound.Scaffolder {
	return &service{}
}

func (s *service) Scaffold(path string, recipe *recipe.Recipe) error {
	fmt.Printf("\n[i] Scaffolding templates for a '%s' project...\n", recipe.Project.Type)

	// 1. Copy generic templates
	if err := s.copyTemplates(filepath.Join("templates", "generic"), path, recipe); err != nil {
		return err
	}

	// 2. Copy stack-specific templates
	if recipe.Project.Type == "go-cli" {
		if err := s.copyTemplates(filepath.Join("templates", "go-cli"), path, recipe); err != nil {
			return err
		}
	}

	if recipe.Persistence.Type == "postgresql" {
		if err := s.copyTemplates(filepath.Join("templates", "postgresql"), path, recipe); err != nil {
			return err
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
		targetPath := filepath.Join(targetDir, templatePath[len(sourceDir):])

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read the template file.
		rawContent, err := templateFiles.ReadFile(templatePath)
		if err != nil {
			return err
		}

		// Execute the template to replace variables like {{ .Project.Name }}
		tmpl, err := template.New(d.Name()).Parse(string(rawContent))
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
