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
}

func NewService(coverageParsers map[string]outbound.CoverageParser, secretScanner outbound.SecretScanner) inbound.ProjectVerifier {
	return &service{
		coverageParsers: coverageParsers,
		secretScanner:   secretScanner,
	}
}

func (s *service) VerifyProject(options inbound.VerifyOptions) error {
	fmt.Println("Running verifications...")

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
			// This is not a fatal error, so we continue.
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
