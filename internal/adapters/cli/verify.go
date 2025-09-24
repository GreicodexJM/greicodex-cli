package cli

import (
	"fmt"
	"grei-cli/internal/adapters/coverage"
	"grei-cli/internal/adapters/linter"
	"grei-cli/internal/adapters/scanner"
	"grei-cli/internal/adapters/syschecker"
	"grei-cli/internal/core/recipe"
	"grei-cli/internal/core/verifier"
	"grei-cli/internal/ports/inbound"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [path]",
	Short: "Verifica que un proyecto existente cumpla con los estándares de Greicodex.",
	Long: `Ejecuta una serie de comprobaciones en un repositorio existente, incluyendo
escaneo de secretos, linters, pruebas, cobertura mínima y configuración de CI/CD,
Helm y OpenTofu.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		// 1. Read and parse the project's recipe file.
		recipePath := filepath.Join(targetPath, "grei.yml")
		recipeData, err := os.ReadFile(recipePath)
		if err != nil {
			color.Red("❌ Error: No se pudo leer el archivo 'grei.yml' en '%s'. Asegúrate de que el proyecto ha sido inicializado.", targetPath)
			os.Exit(1)
		}

		var projRecipe recipe.Recipe
		if err := yaml.Unmarshal(recipeData, &projRecipe); err != nil {
			color.Red("❌ Error: No se pudo parsear el archivo 'grei.yml': %v", err)
			os.Exit(1)
		}

		minCoverage, _ := cmd.Flags().GetInt("min-cov")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// For now, we assume a single coverage parser. This can be expanded later.
		coverageParser := coverage.NewJestParser()
		sysChecker := syschecker.New()
		secretScanner := scanner.NewGitleaksScanner(sysChecker)
		linterDetector := linter.NewFsDetector()
		verifyService := verifier.NewService(coverageParser, secretScanner, linterDetector)

		options := inbound.VerifyOptions{
			Path:        targetPath,
			MinCoverage: minCoverage,
			JSONOutput:  jsonOutput,
			Recipe:      &projRecipe,
		}

		if err := verifyService.VerifyProject(options); err != nil {
			fmt.Fprintf(os.Stderr, "Error durante la verificación: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Verificación completada exitosamente.")
	},
}

func AddVerifyCommand(root *cobra.Command) {
	verifyCmd.Flags().Int("min-cov", 80, "Cobertura de pruebas mínima requerida.")
	verifyCmd.Flags().Bool("json", false, "Muestra la salida en formato JSON.")
	root.AddCommand(verifyCmd)
}
