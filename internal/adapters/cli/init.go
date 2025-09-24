package cli

import (
	"fmt"
	"grei-cli/internal/adapters/downloader"
	"grei-cli/internal/adapters/filesystem"
	"grei-cli/internal/adapters/git"
	"grei-cli/internal/core/initializer"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/core/scaffolder"
	"grei-cli/internal/core/stack"
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

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Inicializa un nuevo proyecto con la estructura estándar de Greicodex.",
	Long: `Crea un nuevo proyecto con una estructura de directorios estándar,
archivos de configuración y plantillas iniciales. Este comando te guiará
a través de una serie de preguntas para configurar el 'grei.yml', el archivo
de receta del proyecto.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		recipePath := filepath.Join(targetPath, "grei.yml")
		if _, err := os.Stat(recipePath); err == nil {
			color.Red("❌ Error: Este directorio ya contiene un proyecto 'grei' (grei.yml encontrado).")
			os.Exit(1)
		}

		fmt.Println("🚀 ¡Bienvenido al inicializador de proyectos de Greicodex!")
		fmt.Println("---------------------------------------------------------")

		answers := recipe.Recipe{}

		// 1. Build dynamic survey from the internal stack registry.
		codeStacks, persistenceStacks, deploymentStacks := categorizeStacks()

		projectQuestions := []*survey.Question{
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "¿Cuál es el nombre del proyecto?", Default: generateProjectName()},
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
					Default: "go-cli",
				},
			},
		}

		if err := survey.Ask(projectQuestions, &answers.Project); err != nil {
			handleSurveyError(err)
		}

		// 2. Handle conditional logic for each stack type
		if answers.Project.Type == "Custom" {
			// Fallback to manual configuration
			if err := survey.Ask(baseStackQuestions, &answers.Stack); err != nil {
				handleSurveyError(err)
			}
		} else {
			// Pre-fill based on the chosen stack
			for _, s := range stack.Registry {
				if s.Name == answers.Project.Type {
					answers.Stack.Language = s.Provides.Language
					answers.Stack.Tooling = s.Provides.Tooling
					answers.Stack.DependencyManagement = s.Provides.DependencyManagement
					answers.Stack.BuildReleaseRun = s.Provides.BuildReleaseRun
					color.Green("✓ Usando la pila de código '%s'.", s.Name)
					break
				}
			}
		}

		// Ask about persistence
		persistenceQuestion := &survey.Question{
			Name: "type",
			Prompt: &survey.Select{
				Message: "¿Qué tipo de pila de persistencia usarás?",
				Options: persistenceStacks,
			},
		}
		if err := survey.Ask([]*survey.Question{persistenceQuestion}, &answers.Persistence); err != nil {
			handleSurveyError(err)
		}

		// Ask about deployment
		deploymentQuestion := &survey.Question{
			Name: "type",
			Prompt: &survey.Select{
				Message: "¿Qué tipo de pila de despliegue usarás?",
				Options: deploymentStacks,
			},
		}
		if err := survey.Ask([]*survey.Question{deploymentQuestion}, &answers.Deployment); err != nil {
			handleSurveyError(err)
		}

		// 3. Generate recipe file
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Creando receta del proyecto (grei.yml)..."
		s.Start()

		yamlData, err := yaml.Marshal(&answers)
		if err != nil {
			s.Stop()
			color.Red("Error al generar el archivo YAML: %v", err)
			os.Exit(1)
		}

		if err := os.WriteFile(recipePath, yamlData, 0644); err != nil {
			s.Stop()
			color.Red("Error al escribir el archivo grei.yml: %v", err)
			os.Exit(1)
		}

		s.Stop()
		color.Green("✅ Receta del proyecto creada exitosamente en '%s'.", recipePath)

		// 4. Scaffold initial templates
		fsRepo := filesystem.NewRepository()
		gitRepo := git.NewRepository()
		downloader := downloader.NewGitDownloader()
		initializerService := initializer.NewService(fsRepo, gitRepo, downloader)
		if err := initializerService.InitializeProject(targetPath, true); err != nil {
			color.Red("❌ Error durante la inicialización: %v", err)
			os.Exit(1)
		}

		scaffolderService := scaffolder.NewService()
		if err := scaffolderService.Scaffold(targetPath, &answers); err != nil {
			color.Red("❌ Error durante el scaffolding: %v", err)
			os.Exit(1)
		}

		fmt.Println("\n🚀 ¡Proyecto inicializado exitosamente!")
	},
}

func handleSurveyError(err error) {
	color.Red("Error durante la encuesta: %v", err)
	os.Exit(1)
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

func categorizeStacks() ([]string, []string, []string) {
	codeStacks := []string{"Custom"}
	persistenceStacks := []string{"Ninguna"}
	deploymentStacks := []string{"Ninguno"}

	for _, s := range stack.Registry {
		switch s.Type {
		case "code":
			codeStacks = append(codeStacks, s.Name)
		case "persistence":
			persistenceStacks = append(persistenceStacks, s.Name)
		case "deployment":
			deploymentStacks = append(deploymentStacks, s.Name)
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

func generateProjectName() string {
	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	con := constellations[rand.Intn(len(constellations))]
	return fmt.Sprintf("%s%s", adj, con)
}

func AddInitCommand(root *cobra.Command) {
	root.AddCommand(initCmd)
}
