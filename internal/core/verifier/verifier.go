package verifier

import (
	"fmt"
	"grei-cli/internal/adapters/scanner"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"os"
	"path/filepath"
)

type service struct {
	coverageParsers map[string]outbound.CoverageParser
	secretScanner   outbound.SecretScanner
	linterDetector  outbound.LinterDetector
}

func NewService(
	coverageParsers map[string]outbound.CoverageParser,
	secretScanner outbound.SecretScanner,
	linterDetector outbound.LinterDetector,
) inbound.ProjectVerifier {
	return &service{
		coverageParsers: coverageParsers,
		secretScanner:   secretScanner,
		linterDetector:  linterDetector,
	}
}

func (s *service) VerifyProject(options inbound.VerifyOptions) error {
	fmt.Println("Running verifications...")

	// Verify project against its recipe
	if options.Recipe != nil {
		fmt.Printf("  [i] Verifying project against recipe for '%s'...\n", options.Recipe.Project.Name)
		if err := s.verifyLinter(options); err != nil {
			return err
		}
		if err := s.verifyPersistence(options); err != nil {
			return err
		}
		if err := s.verifyDeployment(options); err != nil {
			return err
		}
	}

	// Check coverage
	coverage, err := s.findAndParseCoverage(options.Path)
	if err != nil {
		return fmt.Errorf("could not parse coverage file: %w", err)
	}

	fmt.Printf("  [i] Found test coverage: %.2f%%\n", coverage)

	if coverage < float64(options.MinCoverage) {
		return fmt.Errorf("test coverage (%.2f%%) is below the required minimum of %d%%", coverage, options.MinCoverage)
	}

	fmt.Printf("  [✓] Test coverage is sufficient (%.2f%% >= %d%%)\n", coverage, options.MinCoverage)

	// Scan for secrets
	fmt.Println("\nScanning for secrets...")
	secrets, err := s.secretScanner.Scan(options.Path)
	if err != nil {
		if err == scanner.ErrGitleaksNotFound {
			fmt.Println("  [!] gitleaks not found, skipping secret scan.")
		} else {
			return fmt.Errorf("secret scanning failed: %w", err)
		}
	}

	if len(secrets) > 0 {
		fmt.Println("  [✗] Found potential secrets:")
		for _, secret := range secrets {
			fmt.Println(secret)
		}
		return fmt.Errorf("potential secrets found")
	}

	fmt.Println("  [✓] No secrets found.")

	// Check for required files and directories
	fmt.Println("\nChecking for required files and directories...")
	requiredPaths := []string{
		"LICENSE",
		"CONTRIBUTING.md",
		"deploy/helm",
	}

	allPathsExist := true
	for _, p := range requiredPaths {
		if _, err := os.Stat(filepath.Join(options.Path, p)); os.IsNotExist(err) {
			fmt.Printf("  [✗] Missing: %s\n", p)
			allPathsExist = false
		} else {
			fmt.Printf("  [✓] Found: %s\n", p)
		}
	}

	if !allPathsExist {
		return fmt.Errorf("missing required files or directories")
	}

	return nil
}

func (s *service) verifyLinter(options inbound.VerifyOptions) error {
	if options.Recipe.Stack.Linter == "" {
		fmt.Println("  [i] No linter specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying linter '%s'...\n", options.Recipe.Stack.Linter)
	found, err := s.linterDetector.CheckConfig(options.Path, options.Recipe.Stack.Linter)
	if err != nil {
		return fmt.Errorf("could not check for linter config: %w", err)
	}

	if !found {
		return fmt.Errorf("linter config for '%s' not found", options.Recipe.Stack.Linter)
	}

	fmt.Printf("  [✓] Linter config found for '%s'.\n", options.Recipe.Stack.Linter)
	return nil
}

func (s *service) verifyPersistence(options inbound.VerifyOptions) error {
	if options.Recipe.Persistence.Type == "" || options.Recipe.Persistence.Type == "Ninguna" {
		fmt.Println("  [i] No persistence layer specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying persistence layer '%s'...\n", options.Recipe.Persistence.Type)
	// For now, we just check for a docker-compose.yml file.
	// This could be expanded to check for specific migrations, etc.
	composePath := filepath.Join(options.Path, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found for persistence layer '%s'", options.Recipe.Persistence.Type)
	}

	fmt.Printf("  [✓] Found docker-compose.yml for '%s'.\n", options.Recipe.Persistence.Type)
	return nil
}

func (s *service) verifyDeployment(options inbound.VerifyOptions) error {
	if options.Recipe.Deployment.Type == "" || options.Recipe.Deployment.Type == "Ninguno" {
		fmt.Println("  [i] No deployment layer specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying deployment layer '%s'...\n", options.Recipe.Deployment.Type)
	// For now, we just check for a deploy/ directory.
	// This could be expanded to check for specific IaC files, etc.
	deployPath := filepath.Join(options.Path, "deploy")
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return fmt.Errorf("deploy/ directory not found for deployment layer '%s'", options.Recipe.Deployment.Type)
	}

	fmt.Printf("  [✓] Found deploy/ directory for '%s'.\n", options.Recipe.Deployment.Type)
	return nil
}

func (s *service) findAndParseCoverage(basePath string) (float64, error) {
	for name, parser := range s.coverageParsers {
		searchPath := filepath.Join(basePath, "**", name) // Using glob pattern
		matches, err := filepath.Glob(searchPath)
		if err != nil {
			continue // Ignore errors in globbing
		}
		if len(matches) > 0 {
			return parser.Parse(matches[0])
		}
	}
	return 0, os.ErrNotExist
}
