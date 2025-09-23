package doctor

import (
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
)

var requiredTools = []string{
	"git",
	"bash",
	"make",
	"docker",
	"tofu",
	"helm",
	"kubectl",
	"zip",
	"jq",
	"yq",
}

type service struct {
	sysChecker outbound.SystemChecker
}

func NewService(sysChecker outbound.SystemChecker) inbound.DoctorService {
	return &service{
		sysChecker: sysChecker,
	}
}

func (s *service) CheckEnvironment() []inbound.CheckResult {
	results := make([]inbound.CheckResult, 0, len(requiredTools))
	for _, tool := range requiredTools {
		results = append(results, inbound.CheckResult{
			Command: tool,
			Found:   s.sysChecker.CommandExists(tool),
		})
	}
	return results
}
