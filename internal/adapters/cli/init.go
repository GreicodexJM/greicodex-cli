package cli

import (
	"fmt"
	"grei-cli/internal/adapters/downloader"
	"grei-cli/internal/adapters/filesystem"
	"grei-cli/internal/adapters/git"
	"grei-cli/internal/core/initializer"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/core/scaffolder"
	"grei-cli/internal/ports/inbound"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	noInteractive bool
	recipeFile    string
)

// AddInitCommand adds the init command to the root command.
func AddInitCommand(root *cobra.Command) {
	fsRepo := filesystem.NewRepository()
	gitRepo := git.NewRepository()
	initializerService := initializer.NewService(fsRepo, gitRepo)
	scaffolderService := scaffolder.NewService(fsRepo)

	cmd := NewInitCommand(initializerService, scaffolderService)
	cmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Desactiva el modo interactivo y usa un archivo de receta")
	cmd.Flags().StringVar(&recipeFile, "recipe-file", "", "Ruta al archivo de receta (grei.yml) para usar en modo no interactivo")
	root.AddCommand(cmd)
}

// NewInitCommand creates a new init command with its dependencies.
func NewInitCommand(initializerService inbound.InitializerService, scaffolderService inbound.ScaffolderService) *cobra.Command {
	return &cobra.Command{
		Use:   "init [path]",
		Short: "Inicializa un nuevo proyecto con la estructura estándar de Greicodex.",
		Long: `Crea un nuevo proyecto con una estructura de directorios estándar,
archivos de configuración y plantillas iniciales. Este comando te guiará
a través de una serie de preguntas para configurar el 'grei.yml', el archivo
de receta del proyecto.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("error getting user home directory: %w", err)
			}
			cacheDir := filepath.Join(homeDir, ".grei", "templates")

			downloader := downloader.NewGitDownloader()
			if err := downloader.Download(cmd.Context(), "https://github.com/GreicodexJM/greicodex-cli.git", "master", cacheDir); err != nil {
				color.Yellow("Could not download remote templates: %v", err)
			}

			targetPath := "."
			if len(args) > 0 {
				targetPath = args[0]
			}

			recipePath := filepath.Join(targetPath, "grei.yml")
			if _, err := os.Stat(recipePath); err == nil {
				return fmt.Errorf("este directorio ya contiene un proyecto 'grei' (grei.yml encontrado)")
			}

			answers := recipe.Recipe{}

			if noInteractive {
				if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
					return fmt.Errorf("error al crear el directorio del proyecto: %w", err)
				}
				if recipeFile == "" {
					return fmt.Errorf("--recipe-file es requerido en modo no interactivo")
				}
				yamlFile, err := os.ReadFile(recipeFile)
				if err != nil {
					return fmt.Errorf("error al leer el archivo de receta: %w", err)
				}
				err = yaml.Unmarshal(yamlFile, &answers)
				if err != nil {
					return fmt.Errorf("error al parsear el archivo de receta: %w", err)
				}
			} else {
				fmt.Println("🚀 ¡Bienvenido al inicializador de proyectos de Greicodex!")
				fmt.Println("---------------------------------------------------------")

				codeStacks, persistenceStacks, deploymentStacks := CategorizeStacks(cacheDir)

				projectQuestions := []*survey.Question{
					{
						Name:     "name",
						Prompt:   &survey.Input{Message: "¿Cuál es el nombre del proyecto?", Default: GenerateProjectName()},
						Validate: survey.Required,
					},
					{
						Name:     "customer",
						Prompt:   &survey.Input{Message: "¿Quién es el cliente para este proyecto?", Default: "Greicodex"},
						Validate: survey.Required,
					},
					{
						Name: "type",
						Prompt: &survey.Select{
							Message: "¿Qué tipo de pila de código usarás?",
							Options: codeStacks,
							Default: "cli",
						},
					},
				}

				if err := survey.Ask(projectQuestions, &answers.Project); err != nil {
					return fmt.Errorf("error durante la encuesta: %w", err)
				}

				if answers.Project.Type == "Custom" {
					if err := survey.Ask(baseStackQuestions, &answers.Stack); err != nil {
						return fmt.Errorf("error durante la encuesta: %w", err)
					}
				} else {
					// This part will be refactored to use the dynamic template discovery
				}

				persistenceQuestion := &survey.Question{
					Name: "type",
					Prompt: &survey.Select{
						Message: "¿Qué tipo de pila de persistencia usarás?",
						Options: persistenceStacks,
					},
				}
				if err := survey.Ask([]*survey.Question{persistenceQuestion}, &answers.Persistence); err != nil {
					return fmt.Errorf("error durante la encuesta: %w", err)
				}

				deploymentQuestion := &survey.Question{
					Name: "type",
					Prompt: &survey.Select{
						Message: "¿Qué tipo de pila de despliegue usarás?",
						Options: deploymentStacks,
					},
				}
				if err := survey.Ask([]*survey.Question{deploymentQuestion}, &answers.Deployment); err != nil {
					return fmt.Errorf("error durante la encuesta: %w", err)
				}
			}

			s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " Creando receta del proyecto (grei.yml)..."
			s.Start()

			yamlData, err := yaml.Marshal(&answers)
			if err != nil {
				s.Stop()
				return fmt.Errorf("error al generar el archivo YAML: %w", err)
			}

			if err := os.WriteFile(recipePath, yamlData, 0644); err != nil {
				s.Stop()
				return fmt.Errorf("error al escribir el archivo grei.yml: %w", err)
			}

			s.Stop()
			color.Green("✅ Receta del proyecto creada exitosamente en '%s'.", recipePath)

			if err := initializerService.InitializeProject(targetPath, cacheDir, true); err != nil {
				return fmt.Errorf("error durante la inicialización: %w", err)
			}

			if err := scaffolderService.Scaffold(targetPath, &answers); err != nil {
				return fmt.Errorf("error durante el scaffolding: %w", err)
			}

			fmt.Println("\n🚀 ¡Proyecto inicializado exitosamente!")
			return nil
		},
	}
}

var baseStackQuestions = []*survey.Question{
	{
		Name:   "language",
		Prompt: &survey.Select{Message: "¿Qué lenguaje principal usarás?", Options: []string{"Go", "TypeScript", "Python"}, Default: "Go"},
	},
	{
		Name:   "tooling",
		Prompt: &survey.Input{Message: "¿Qué tooling principal (framework, etc.) usarás?"},
	},
}

func CategorizeStacks(cacheDir string) ([]string, []string, []string) {
	codeStacks := []string{"Custom"}
	persistenceStacks := []string{"Ninguna"}
	deploymentStacks := []string{"Ninguno"}

	templateDirs, err := os.ReadDir(filepath.Join(cacheDir, "templates"))
	if err != nil {
		return codeStacks, persistenceStacks, deploymentStacks
	}

	for _, dir := range templateDirs {
		if !dir.IsDir() || dir.Name() == "generic" {
			continue
		}

		manifestPath := filepath.Join(cacheDir, "templates", dir.Name(), "manifest.yml")
		manifestFile, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var manifest scaffolder.Manifest
		if err := yaml.Unmarshal(manifestFile, &manifest); err != nil {
			continue
		}

		switch manifest.Type {
		case "code":
			codeStacks = append(codeStacks, manifest.Name)
		case "persistence":
			persistenceStacks = append(persistenceStacks, manifest.Name)
		case "deployment":
			deploymentStacks = append(deploymentStacks, manifest.Name)
		}
	}

	return codeStacks, persistenceStacks, deploymentStacks
}

var adjectives = []string{
	"Adaptable", "Agil", "Alegre", "Ambicioso", "Amable", "Audaz", "Brillante", "Calmado", "Capaz", "Carismatico",
	"Compasivo", "Confiable", "Creativo", "Curioso", "Decidido", "Dedicado", "Dinamico", "Eficiente", "Elegante", "Empatico",
	"Energico", "Entusiasta", "Estelar", "Exitoso", "Flexible", "Fuerte", "Honesto", "Imaginativo", "Innovador", "Inspirador",
	"Inteligente", "Intrepido", "Jovial", "Leal", "Luminoso", "Metodico", "Optimista", "Organizado", "Paciente", "Perseverante",
	"Poderoso", "Positivo", "Practico", "Radiante", "Resiliente", "Sabio", "Seguro", "Tenaz", "Valiente", "Visionario",
}

var constellations = []string{
	"Andromeda", "Antlia", "Apus", "Aquarius", "Aquila", "Ara", "Aries", "Auriga", "Bootes", "Caelum",
	"Camelopardalis", "Cancer", "CanesVenatici", "CanisMajor", "CanisMinor", "Capricornus", "Carina", "Cassiopeia", "Centaurus", "Cepheus",
	"Cetus", "Chamaeleon", "Circinus", "Columba", "ComaBerenices", "CoronaAustralis", "CoronaBorealis", "Corvus", "Crater", "Crux",
	"Cygnus", "Delphinus", "Dorado", "Draco", "Equuleus", "Eridanus", "Fornax", "Gemini", "Grus", "Hercules",
	"Horologium", "Hydra", "Hydrus", "Indus", "Lacerta", "Leo", "LeoMinor", "Lepus", "Libra", "Lupus",
	"Lynx", "Lyra", "Mensa", "Microscopium", "Monoceros", "Musca", "Norma", "Octans", "Ophiuchus", "Orion",
	"Pavo", "Pegasus", "Perseus", "Phoenix", "Pictor", "Pisces", "PiscisAustrinus", "Puppis", "Pyxis", "Reticulum",
	"Sagitta", "Sagittarius", "Scorpius", "Sculptor", "Scutum", "Serpens", "Sextans", "Taurus", "Telescopium", "Triangulum",
	"TriangulumAustrale", "Tucana", "UrsaMajor", "UrsaMinor", "Vela", "Virgo", "Volans", "Vulpecula",
}

func GenerateProjectName() string {
	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	con := constellations[rand.Intn(len(constellations))]
	return fmt.Sprintf("%s%s", adj, con)
}
