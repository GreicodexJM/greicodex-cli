package coverage

import (
	"encoding/json"
	"grei-cli/internal/ports/outbound"
	"os"
)

type JestCoverageParser struct{}

func NewJestParser() outbound.CoverageParser {
	return &JestCoverageParser{}
}

type jestCoverageSummary struct {
	Total struct {
		Lines struct {
			Pct float64 `json:"pct"`
		} `json:"lines"`
	} `json:"total"`
}

func (p *JestCoverageParser) Parse(path string) (float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	var summary jestCoverageSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return 0, err
	}

	return summary.Total.Lines.Pct, nil
}
