package scanner

import (
	"fmt"
	"grei-cli/internal/ports/outbound"
	"os/exec"
)

type GitleaksScanner struct {
	sysChecker outbound.SystemChecker
}

func NewGitleaksScanner(sysChecker outbound.SystemChecker) outbound.SecretScanner {
	return &GitleaksScanner{
		sysChecker: sysChecker,
	}
}

var ErrGitleaksNotFound = fmt.Errorf("gitleaks not found")

func (s *GitleaksScanner) Scan(path string) ([]string, error) {
	if !s.sysChecker.CommandExists("gitleaks") {
		return nil, ErrGitleaksNotFound
	}

	cmd := exec.Command("gitleaks", "detect", "-s", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// gitleaks exits with a non-zero status code if it finds secrets.
		// We need to check the output to see if secrets were found.
		if len(output) > 0 {
			return []string{string(output)}, nil
		}
		return nil, fmt.Errorf("gitleaks failed: %w\n%s", err, output)
	}

	return nil, nil // No secrets found
}
