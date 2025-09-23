package initializer

import (
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"grei-cli/internal/templates"
	"path/filepath"
	"time"
)

type service struct {
	fsRepo  outbound.FSRepository
	gitRepo outbound.GitRepository
}

func NewService(fsRepo outbound.FSRepository, gitRepo outbound.GitRepository) inbound.ProjectInitializer {
	return &service{
		fsRepo:  fsRepo,
		gitRepo: gitRepo,
	}
}

func (s *service) InitializeProject(path string, gitInit bool) error {
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
		content, err := templates.Process(tmpl, data)
		if err != nil {
			return err
		}
		err = s.fsRepo.CreateFile(filepath.Join(path, dest), content)
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
