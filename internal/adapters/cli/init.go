package cli

import (
	"fmt"
	"grei-cli/internal/core/recipe"
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

// surveyAnswers is a flat struct to avoid issues with survey's nested struct mapping.
type surveyAnswers struct {
	ProjectName        string
	ProjectType        string
	ProjectLanguage    string
	StackFramework     string
	StackLinter        string
	StackTesting       string
	StackCICD          []string
	DeploymentType     string
	DeploymentProvider string
}

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
		fmt.Println("Por favor, responde las siguientes preguntas para crear tu receta de proyecto.")

		flatAnswers := surveyAnswers{}
		err := survey.Ask(questions, &flatAnswers)
		if err != nil {
			color.Red("Error durante la encuesta: %v", err)
			os.Exit(1)
		}

		// Manually populate the nested recipe struct from the flat survey answers.
		answers := recipe.Recipe{
			Project: recipe.Project{
				Name:     flatAnswers.ProjectName,
				Type:     flatAnswers.ProjectType,
				Language: flatAnswers.ProjectLanguage,
			},
			Stack: recipe.Stack{
				Framework: flatAnswers.StackFramework,
				Linter:    flatAnswers.StackLinter,
				Testing:   flatAnswers.StackTesting,
				CICD:      flatAnswers.StackCICD,
			},
			Deployment: recipe.Deployment{
				Type:     flatAnswers.DeploymentType,
				Provider: flatAnswers.DeploymentProvider,
			},
		}

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = " Creando receta del proyecto (grei.yml)..."
		s.Start()

		yamlData, err := yaml.Marshal(&answers)
		if err != nil {
			s.Stop()
			color.Red("Error al generar el archivo YAML: %v", err)
			os.Exit(1)
		}

		err = os.WriteFile(recipePath, yamlData, 0644)
		if err != nil {
			s.Stop()
			color.Red("Error al escribir el archivo grei.yml: %v", err)
			os.Exit(1)
		}

		s.Stop()
		color.Green("✅ Receta del proyecto creada exitosamente en '%s'.", recipePath)

		scaffoldTemplates(targetPath, &answers)

		fmt.Println("\n🚀 ¡Proyecto inicializado exitosamente!")
	},
}

var questions = []*survey.Question{
	{
		Name:     "ProjectName",
		Prompt:   &survey.Input{Message: "¿Cuál es el nombre del proyecto?", Default: generateProjectName()},
		Validate: survey.Required,
	},
	{
		Name:     "CustomerName",
		Prompt:   &survey.Input{Message: "¿Cuál es el nombre del cliente?", Default: "Greicodex"},
		Validate: survey.Required,
	},
	{
		Name: "ProjectType",
		Prompt: &survey.Select{Message: "¿Qué tipo de proyecto es?", Options: []string{
			"Interfaz Movil",
			"Interfaz Web",
			"Interfaz Desktop",
			"API de Servicio",
			"API Serverless",
			"Proceso por Lotes",
			"Tarea Programada",
			"Herramienta UI",
			"Herramienta CLI",
		}, Default: "Servicio Backend"},
	},
	{
		Name:   "ProjectLanguage",
		Prompt: &survey.Select{Message: "¿Qué lenguaje principal usarás?", Options: []string{"Go", "TypeScript", "Python"}, Default: "Go"},
	},
	{
		Name:   "StackFramework",
		Prompt: &survey.Select{Message: "¿Qué framework principal usarás?", Options: []string{"Gin (Go)", "Cobra (Go)", "React (TypeScript)", "FastAPI (Python)", "Ninguno"}, Default: "Cobra (Go)"},
	},
	{
		Name:   "StackLinter",
		Prompt: &survey.Select{Message: "¿Qué linter prefieres?", Options: []string{"golangci-lint (Go)", "ESLint (TypeScript)", "Ruff (Python)"}, Default: "golangci-lint (Go)"},
	},
	{
		Name:   "StackTesting",
		Prompt: &survey.Select{Message: "¿Qué herramienta de pruebas usarás?", Options: []string{"go-test (Go)", "Jest (TypeScript)", "Pytest (Python)"}, Default: "go-test (Go)"},
	},
	{
		Name:   "StackCICD",
		Prompt: &survey.MultiSelect{Message: "¿Qué herramientas de CI/CD necesitas?", Options: []string{"GitHub Actions", "GitLab CI", "CircleCI"}},
	},
	{
		Name:   "DeploymentType",
		Prompt: &survey.Select{Message: "¿Cuál es el tipo de despliegue?", Options: []string{"Kubernetes", "Serverless", "Docker Swarm", "Binario"}, Default: "Kubernetes"},
	},
	{
		Name:   "DeploymentProvider",
		Prompt: &survey.Select{Message: "¿Cuál es el proveedor de despliegue?", Options: []string{"GCP", "AWS", "Azure", "On-premise"}, Default: "GCP"},
	},
}

func scaffoldTemplates(path string, recipe *recipe.Recipe) {
	fmt.Printf("\n[i] Scaffolding templates for a '%s' project...\n", recipe.Project.Language)
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
	return fmt.Sprintf("%s-%s", con, adj)
}

func AddInitCommand(root *cobra.Command) {
	root.AddCommand(initCmd)
}
