package scaffolder

import (
	"encoding/json"
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
	Path     string        `json:"path"`
	Strategy MergeStrategy `json:"strategy"`
}

type SkeletonManifest struct {
	Files []FileManifest `json:"files"`
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
		manifestPath := filepath.Join(skeleton, "manifest.json")
		manifestData, err := s.fsRepo.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("failed to read manifest file %s: %w", manifestPath, err)
		}

		var manifest SkeletonManifest
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			return fmt.Errorf("failed to unmarshal manifest file %s: %w", manifestPath, err)
		}

		for _, file := range manifest.Files {
			sourcePath := filepath.Join(skeleton, file.Path)
			targetPath := filepath.Join(targetDir, file.Path)

			sourceContent, err := s.fsRepo.ReadFile(sourcePath)
			if err != nil {
				return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
			}

			if err := s.applyStrategy(file.Strategy, sourceContent, targetPath); err != nil {
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
