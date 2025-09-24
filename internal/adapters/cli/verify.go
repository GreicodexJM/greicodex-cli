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

// AddVerifyCommand adds the verify command to the root command.
func AddVerifyCommand(root *cobra.Command) {
	coverageParser := coverage.NewGoParser()
	sysChecker := syschecker.New()
	secretScanner := scanner.NewGitleaksScanner(sysChecker)
	linterDetector := linter.NewFsDetector()
	verifyService := verifier.NewService(coverageParser, secretScanner, linterDetector)

	cmd := NewVerifyCommand(verifyService)
	cmd.Flags().Int("min-cov", 80, "Cobertura de pruebas mínima requerida.")
	cmd.Flags().Bool("json", false, "Muestra la salida en formato JSON.")
	root.AddCommand(cmd)
}

// NewVerifyCommand creates a new verify command with its dependencies.
func NewVerifyCommand(verifyService inbound.VerifierService) *cobra.Command {
	return &cobra.Command{
		Use:   "verify [path]",
		Short: "Verifica que un proyecto existente cumpla con los estándares de Greicodex.",
		Long: `Ejecuta una serie de comprobaciones en un repositorio existente, incluyendo
escaneo de secretos, linters, pruebas, cobertura mínima y configuración de CI/CD,
Helm y OpenTofu.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := "."
			if len(args) > 0 {
				targetPath = args[0]
			}

			recipePath := filepath.Join(targetPath, "grei.yml")
			recipeData, err := os.ReadFile(recipePath)
			if err != nil {
				return fmt.Errorf("no se pudo leer el archivo 'grei.yml' en '%s'. Asegúrate de que el proyecto ha sido inicializado", targetPath)
			}

			var projRecipe recipe.Recipe
			if err := yaml.Unmarshal(recipeData, &projRecipe); err != nil {
				return fmt.Errorf("no se pudo parsear el archivo 'grei.yml': %w", err)
			}

			minCoverage, _ := cmd.Flags().GetInt("min-cov")
			jsonOutput, _ := cmd.Flags().GetBool("json")

			options := inbound.VerifyOptions{
				Path:        targetPath,
				MinCoverage: minCoverage,
				JSONOutput:  jsonOutput,
				Recipe:      &projRecipe,
			}

			if err := verifyService.VerifyProject(options); err != nil {
				return fmt.Errorf("error durante la verificación: %w", err)
			}

			color.Green("¡Proyecto verificado exitosamente!")
			return nil
		},
	}
}
