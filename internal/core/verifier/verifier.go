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
	coverageParser outbound.CoverageParser
	secretScanner  outbound.SecretScanner
	linterDetector outbound.LinterDetector
}

func NewService(
	coverageParser outbound.CoverageParser,
	secretScanner outbound.SecretScanner,
	linterDetector outbound.LinterDetector,
) inbound.VerifierService {
	return &service{
		coverageParser: coverageParser,
		secretScanner:  secretScanner,
		linterDetector: linterDetector,
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
	linter, ok := options.Recipe.Stack["linter"].(string)
	if !ok || linter == "" {
		fmt.Println("  [i] No linter specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying linter '%s'...\n", linter)
	found, err := s.linterDetector.CheckConfig(options.Path, linter)
	if err != nil {
		return fmt.Errorf("could not check for linter config: %w", err)
	}

	if !found {
		return fmt.Errorf("linter config for '%s' not found", linter)
	}

	fmt.Printf("  [✓] Linter config found for '%s'.\n", linter)
	return nil
}

func (s *service) verifyPersistence(options inbound.VerifyOptions) error {
	persistence, ok := options.Recipe.Stack["persistence"].(string)
	if !ok || persistence == "" || persistence == "None" {
		fmt.Println("  [i] No persistence layer specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying persistence layer '%s'...\n", persistence)
	// For now, we just check for a docker-compose.yml file.
	// This could be expanded to check for specific migrations, etc.
	composePath := filepath.Join(options.Path, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found for persistence layer '%s'", persistence)
	}

	fmt.Printf("  [✓] Found docker-compose.yml for '%s'.\n", persistence)
	return nil
}

func (s *service) verifyDeployment(options inbound.VerifyOptions) error {
	deployment, ok := options.Recipe.Stack["deployment"].(string)
	if !ok || deployment == "" || deployment == "None" {
		fmt.Println("  [i] No deployment layer specified in recipe, skipping check.")
		return nil
	}

	fmt.Printf("  [i] Verifying deployment layer '%s'...\n", deployment)
	// For now, we just check for a deploy/ directory.
	// This could be expanded to check for specific IaC files, etc.
	deployPath := filepath.Join(options.Path, "deploy")
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return fmt.Errorf("deploy/ directory not found for deployment layer '%s'", deployment)
	}

	fmt.Printf("  [✓] Found deploy/ directory for '%s'.\n", deployment)
	return nil
}

func (s *service) findAndParseCoverage(basePath string) (float64, error) {
	// For now, we assume a single coverage file format. This can be expanded later.
	searchPath := filepath.Join(basePath, "coverage.out")
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		return 0, os.ErrNotExist
	}
	return s.coverageParser.Parse(searchPath)
}
