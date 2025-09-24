package stack

// Registry holds all the available built-in stacks.
var Registry = []*Stack{
	{
		Name:        "symfony-lamp",
		Description: "Pila tecnológica LAMP para aplicaciones web y APIs con Symfony.",
		Type:        "code",
		Provides: Provides{
			Language:             "PHP",
			Tooling:              "Symfony",
			Runtime:              "Apache",
			Persistence:          "MySQL",
			DependencyManagement: "Composer",
			BuildReleaseRun:      "Standard PHP build process",
		},
	},
	{
		Name:        "mern",
		Description: "Pila tecnológica MERN para aplicaciones web full-stack.",
		Type:        "code",
		Provides: Provides{
			Language:             "TypeScript",
			Tooling:              "React",
			Runtime:              "Node.js",
			Persistence:          "MongoDB",
			DependencyManagement: "NPM",
			BuildReleaseRun:      "npm build, npm start",
		},
	},
	{
		Name:        "go-cli",
		Description: "Pila tecnológica para herramientas CLI en Go con Cobra.",
		Type:        "code",
		Provides: Provides{
			Language:             "Go",
			Tooling:              "Cobra",
			Runtime:              "Binario",
			Persistence:          "Filesystem",
			DependencyManagement: "Go Modules",
			BuildReleaseRun:      "go build, ./binary",
		},
	},
	{
		Name:        "postgresql",
		Description: "Stack de persistencia para PostgreSQL.",
		Type:        "persistence",
		Provides: Provides{
			Persistence: "PostgreSQL",
		},
	},
	{
		Name:        "kubernetes",
		Description: "Stack de despliegue para Kubernetes.",
		Type:        "deployment",
		Provides: Provides{
			Runtime: "Kubernetes",
		},
	},
	{
		Name:        "serverless",
		Description: "Stack de despliegue para Serverless.",
		Type:        "deployment",
		Provides: Provides{
			Runtime: "Serverless",
		},
	},
}
