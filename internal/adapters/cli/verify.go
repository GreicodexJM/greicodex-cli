package cli

import (
	"fmt"
	"grei-cli/internal/adapters/coverage"
	"grei-cli/internal/adapters/scanner"
	"grei-cli/internal/adapters/syschecker"
	"grei-cli/internal/core/verifier"
	"grei-cli/internal/ports/inbound"
	"grei-cli/internal/ports/outbound"
	"os"

	"github.com/spf13/cobra"
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

		minCoverage, _ := cmd.Flags().GetInt("min-cov")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		coverageParsers := map[string]outbound.CoverageParser{
			"coverage-summary.json": coverage.NewJestParser(),
		}
		sysChecker := syschecker.New()
		secretScanner := scanner.NewGitleaksScanner(sysChecker)
		verifyService := verifier.NewService(coverageParsers, secretScanner)

		options := inbound.VerifyOptions{
			Path:        targetPath,
			MinCoverage: minCoverage,
			JSONOutput:  jsonOutput,
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
