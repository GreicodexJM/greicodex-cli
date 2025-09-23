package inbound

type CheckResult struct {
	Command string
	Found   bool
}

// DoctorService defines the port for the environment checking service.
type DoctorService interface {
	CheckEnvironment() []CheckResult
}
