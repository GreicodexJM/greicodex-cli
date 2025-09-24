package cli

import (
	"fmt"
	"grei-cli/internal/adapters/plugin"
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

		// 1. Discover available plugins
		pluginScanner := plugin.NewScanner()
		manifests, err := pluginScanner.Scan()
		if err != nil {
			color.Red("Error al escanear plugins: %v", err)
			os.Exit(1)
		}

		fmt.Println("🚀 ¡Bienvenido al inicializador de proyectos de Greicodex!")
		fmt.Println("---------------------------------------------------------")

		answers := recipe.Recipe{}

		// 2. Build dynamic survey based on plugins, ensuring no duplicates.
		projectTypeSet := make(map[string]bool)
		for _, m := range manifests {
			for _, v := range m.Provides {
				projectTypeSet[v.ProjectType] = true
			}
		}

		projectTypes := []string{"Custom"}
		for pt := range projectTypeSet {
			projectTypes = append(projectTypes, pt)
		}

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
					Message: "¿Qué tipo de aplicación estás construyendo?",
					Options: projectTypes,
					Default: "API / Backend",
				},
			},
		}

		if err := survey.Ask(projectQuestions, &answers.Project); err != nil {
			handleSurveyError(err)
		}

		// 3. Handle conditional logic
		if answers.Project.Type == "Custom" {
			// Fallback to manual configuration if no plugin is chosen
			if err := survey.Ask(baseStackQuestions, &answers.Stack); err != nil {
				handleSurveyError(err)
			}
		} else {
			// Pre-fill based on the chosen plugin vertical
			for _, m := range manifests {
				for _, v := range m.Provides {
					if v.ProjectType == answers.Project.Type {
						// Fully populate the stack from the plugin's vertical
						answers.Stack.Language = v.Stack.Language
						// This is a simplified mapping. A real implementation would be more robust.
						if v.Stack.Framework != "" {
							if answers.Project.Type == "Aplicación Web" {
								if answers.WebApp == nil {
									answers.WebApp = &recipe.WebApp{}
								}
								answers.WebApp.Framework = v.Stack.Framework
							} else if answers.Project.Type == "API / Backend" {
								if answers.Api == nil {
									answers.Api = &recipe.Api{}
								}
								answers.Api.Framework = v.Stack.Framework
							}
						}
						color.Green("✓ Usando la pila tecnológica del plugin '%s'.", m.Name)
						break
					}
				}
			}
		}

		if err := survey.Ask(deploymentQuestions, &answers.Deployment); err != nil {
			handleSurveyError(err)
		}

		// 4. Generate recipe file
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

		scaffoldTemplates(targetPath, &answers)

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
}

var deploymentQuestions = []*survey.Question{
	{
		Name:   "type",
		Prompt: &survey.Select{Message: "¿Cuál es el tipo de despliegue?", Options: []string{"Kubernetes", "Serverless", "Docker Swarm", "Binario"}, Default: "Kubernetes"},
	},
	{
		Name:   "provider",
		Prompt: &survey.Select{Message: "¿Cuál es el proveedor de despliegue?", Options: []string{"GCP", "AWS", "Azure", "On-premise"}, Default: "GCP"},
	},
}

func scaffoldTemplates(path string, recipe *recipe.Recipe) {
	fmt.Printf("\n[i] Scaffolding templates for a '%s' project...\n", recipe.Project.Type)
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
