package scaffolder

import (
	"fmt"
	"grei-cli/internal/ports/outbound"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MergeStrategy string

const (
	Overwrite MergeStrategy = "overwrite"
	MergeYAML MergeStrategy = "merge-yaml"
	Append    MergeStrategy = "append"
	Skip      MergeStrategy = "skip"
)

type FileManifest struct {
	Path     string        `yaml:"path"`
	Strategy MergeStrategy `yaml:"strategy"`
}

type MergerService struct {
	fsRepo outbound.FSRepository
}

func NewMergerService(fsRepo outbound.FSRepository) *MergerService {
	return &MergerService{
		fsRepo: fsRepo,
	}
}

func (s *MergerService) Merge(skeletons []string, targetDir string) error {
	for _, skeleton := range skeletons {
		manifestPath := filepath.Join(skeleton, "manifest.yml")
		manifestData, err := s.fsRepo.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("failed to read manifest file %s: %w", manifestPath, err)
		}

		var manifest Manifest
		if err := yaml.Unmarshal(manifestData, &manifest); err != nil {
			return fmt.Errorf("failed to unmarshal manifest file %s: %w", manifestPath, err)
		}

		files, err := s.fsRepo.ReadDir(skeleton)
		if err != nil {
			return fmt.Errorf("failed to read skeleton directory %s: %w", skeleton, err)
		}

		for _, file := range files {
			if file.IsDir() || file.Name() == "manifest.yml" {
				continue
			}

			sourcePath := filepath.Join(skeleton, file.Name())
			targetPath := filepath.Join(targetDir, file.Name())

			sourceContent, err := s.fsRepo.ReadFile(sourcePath)
			if err != nil {
				return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
			}

			strategy := Skip
			for _, f := range manifest.Files {
				if f.Path == file.Name() {
					strategy = f.Strategy
					break
				}
			}

			if err := s.applyStrategy(strategy, sourceContent, targetPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MergerService) applyStrategy(strategy MergeStrategy, sourceContent []byte, targetPath string) error {
	_, err := os.Stat(targetPath)
	fileExists := !os.IsNotExist(err)

	switch strategy {
	case Overwrite:
		return s.fsRepo.CreateFile(targetPath, sourceContent)
	case MergeYAML:
		if !fileExists {
			return s.fsRepo.CreateFile(targetPath, sourceContent)
		}
		return s.mergeYAML(sourceContent, targetPath)
	case Append:
		if !fileExists {
			return s.fsRepo.CreateFile(targetPath, sourceContent)
		}
		return s.appendToFile(sourceContent, targetPath)
	case Skip:
		if fileExists {
			return nil
		}
		return s.fsRepo.CreateFile(targetPath, sourceContent)
	default:
		return fmt.Errorf("unknown merge strategy: %s", strategy)
	}
}

func (s *MergerService) mergeYAML(sourceContent []byte, targetPath string) error {
	targetContent, err := s.fsRepo.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target file for merging: %w", err)
	}

	var sourceData, targetData map[string]interface{}
	if err := yaml.Unmarshal(sourceContent, &sourceData); err != nil {
		return fmt.Errorf("failed to unmarshal source YAML: %w", err)
	}
	if err := yaml.Unmarshal(targetContent, &targetData); err != nil {
		return fmt.Errorf("failed to unmarshal target YAML: %w", err)
	}

	mergedData := mergeMaps(targetData, sourceData)

	mergedContent, err := yaml.Marshal(mergedData)
	if err != nil {
		return fmt.Errorf("failed to marshal merged YAML: %w", err)
	}

	return s.fsRepo.CreateFile(targetPath, mergedContent)
}

func (s *MergerService) appendToFile(sourceContent []byte, targetPath string) error {
	targetContent, err := s.fsRepo.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target file for appending: %w", err)
	}

	newContent := append(targetContent, sourceContent...)
	return s.fsRepo.CreateFile(targetPath, newContent)
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
