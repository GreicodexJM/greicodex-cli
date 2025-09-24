package recipe

// Recipe represents the structure of the grei.yml file.
type Recipe struct {
	Project    Project    `yaml:"project" survey:"project"`
	Stack      Stack      `yaml:"stack,omitempty" survey:"stack"`
	WebApp     *WebApp    `yaml:"webapp,omitempty" survey:"webapp"`
	Api        *Api       `yaml:"api,omitempty" survey:"api"`
	Deployment Deployment `yaml:"deployment,omitempty" survey:"deployment"`
}

// Project contains basic information about the project.
type Project struct {
	Name     string `yaml:"name" survey:"name"`
	Customer string `yaml:"customer" survey:"customer"`
	Type     string `yaml:"type" survey:"type"`
}

// Stack defines the base technology stack used in the project.
type Stack struct {
	Language string   `yaml:"language" survey:"language"`
	Linter   string   `yaml:"linter" survey:"linter"`
	Testing  string   `yaml:"testing" survey:"testing"`
	CICD     []string `yaml:"cicd" survey:"cicd"`
}

// WebApp defines web-specific technologies.
type WebApp struct {
	Framework string `yaml:"framework" survey:"framework"`
}

// Api defines API-specific technologies.
type Api struct {
	Framework string `yaml:"framework" survey:"framework"`
}

// Deployment specifies the deployment target and provider.
type Deployment struct {
	Type     string `yaml:"type" survey:"type"`
	Provider string `yaml:"provider" survey:"provider"`
}
