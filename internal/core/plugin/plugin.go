package plugin

// Manifest is the top-level structure for a plugin's discovery metadata.
type Manifest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Provides    []TechVertical `json:"provides"`
}

// TechVertical represents a complete technology stack for a specific project type.
type TechVertical struct {
	ProjectType string `json:"projectType"`
	Stack       Stack  `json:"stack"`
}

// Stack defines the specific technologies offered by a plugin's vertical.
type Stack struct {
	Language  string `json:"language,omitempty"`
	Framework string `json:"framework,omitempty"`
	WebServer string `json:"webServer,omitempty"`
	Database  string `json:"database,omitempty"`
}
