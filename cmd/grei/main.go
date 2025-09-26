package main

import (
	"fmt"
	"grei-cli/internal/adapters/cli"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grei",
	Short: "GRX CLI es una herramienta para estandarizar proyectos de software.",
	Long: `GRX CLI (Greicodex CLI) es una herramienta de línea de comandos
que automatiza la inicialización, verificación y estandarización de 
proyectos de software bajo los lineamientos de Greicodex.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is provided
		cmd.Help()
	},
}

func init() {
	cli.AddInitCommand(rootCmd)
	cli.AddVerifyCommand(rootCmd)
	cli.AddInstallHooksCommand(rootCmd)
	cli.AddDoctorCommand(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Hubo un error al ejecutar GRX CLI: '%s'\n\n", err)
		os.Exit(1)
	}
}
