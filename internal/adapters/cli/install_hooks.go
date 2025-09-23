package cli

import (
	"fmt"
	"grei-cli/internal/adapters/git"
	"grei-cli/internal/core/hooks"
	"os"

	"github.com/spf13/cobra"
)

var installHooksCmd = &cobra.Command{
	Use:   "install-hooks [path]",
	Short: "Instala los Git hooks estándar de Greicodex en el repositorio.",
	Long: `Configura Git para usar el directorio '.githooks' del repositorio,
asegurando que todos los desarrolladores usen los mismos hooks.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		gitRepo := git.NewRepository()
		hooksService := hooks.NewService(gitRepo)

		if err := hooksService.InstallHooks(targetPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error al instalar los Git hooks: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Git hooks instalados exitosamente.")
	},
}

func AddInstallHooksCommand(root *cobra.Command) {
	root.AddCommand(installHooksCmd)
}
