package recipe

// Recipe represents the structure of the grei.yml file.
type Recipe struct {
	Project    Project    `yaml:"project" survey:"project"`
	Stack      Stack      `yaml:"stack" survey:"stack"`
	Deployment Deployment `yaml:"deployment" survey:"deployment"`
}

// Project contains basic information about the project.
type Project struct {
	Name     string `yaml:"name" survey:"name"`
	Type     string `yaml:"type" survey:"type"`
	Language string `yaml:"language" survey:"language"`
}

// Stack defines the technology stack used in the project.
type Stack struct {
	Framework string   `yaml:"framework" survey:"framework"`
	Linter    string   `yaml:"linter" survey:"linter"`
	Testing   string   `yaml:"testing" survey:"testing"`
	CICD      []string `yaml:"cicd" survey:"cicd"`
}

// Deployment specifies the deployment target and provider.
type Deployment struct {
	Type     string `yaml:"type" survey:"type"`
	Provider string `yaml:"provider" survey:"provider"`
}
