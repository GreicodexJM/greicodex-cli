package coverage

import (
	"bufio"
	"grei-cli/internal/ports/outbound"
	"os"
	"strings"
)

type GoCoverageParser struct{}

func NewGoParser() outbound.CoverageParser {
	return &GoCoverageParser{}
}

func (p *GoCoverageParser) Parse(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var totalStatements, coveredStatements int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			totalStatements += 1
			if parts[2] == "1" {
				coveredStatements += 1
			}
		}
	}

	if totalStatements == 0 {
		return 0, nil
	}

	return (float64(coveredStatements) / float64(totalStatements)) * 100, nil
}
