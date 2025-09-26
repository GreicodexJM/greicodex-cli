package recipe

// Recipe represents the structure of the grei.yml file.
type Recipe struct {
	Project Project                `yaml:"project" survey:"project"`
	Stack   map[string]interface{} `yaml:"stack,omitempty" survey:"stack"`
}

// Project contains basic information about the project.
type Project struct {
	Name     string `yaml:"name" survey:"name"`
	Customer string `yaml:"customer" survey:"customer"`
	Type     string `yaml:"type" survey:"type"`
}
