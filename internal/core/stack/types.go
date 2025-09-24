package stack

// Stack represents a technology stack that can be used in a project.
type Stack struct {
	Name        string
	Description string
	Type        string
	Provides    Provides
}

// Provides defines the technologies that a stack provides.
type Provides struct {
	Language             string
	Tooling              string
	Runtime              string
	Persistence          string
	DependencyManagement string `yaml:"dependencyManagement,omitempty"`
	BuildReleaseRun      string `yaml:"buildReleaseRun,omitempty"`
}
