package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	Language    string `json:"language,omitempty"`
	Framework   string `json:"framework,omitempty"`
	AppServer   string `json:"appServer,omitempty"`
	Persistence string `json:"persistence,omitempty"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "discover" {
		discover()
		return
	}

	// In the future, this will handle the "scaffold" command.
	fmt.Println("Mock Symfony Plugin: use 'discover' command.")
}

func discover() {
	manifest := Manifest{
		Name:        "external-symfony-lamp",
		Description: "Pila tecnológica LAMP externa para Symfony.",
		Provides: []TechVertical{
			{
				ProjectType: "Aplicación Web",
				Stack: Stack{
					Language:  "PHP",
					Framework: "Symfony",
					WebServer: "Apache",
					Database:  "MySQL",
				},
			},
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manifest); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding manifest: %v\n", err)
		os.Exit(1)
	}
}
