package cli

import (
	"fmt"
	"grei-cli/internal/adapters/git"
	"grei-cli/internal/core/hooks"
	"grei-cli/internal/ports/inbound"

	"github.com/spf13/cobra"
)

func AddInstallHooksCommand(root *cobra.Command) {
	gitRepo := git.NewRepository()
	hooksService := hooks.NewService(gitRepo)

	cmd := NewInstallHooksCommand(hooksService)
	root.AddCommand(cmd)
}

func NewInstallHooksCommand(hooksService inbound.HooksService) *cobra.Command {
	return &cobra.Command{
		Use:   "install-hooks [path]",
		Short: "Instala los Git hooks estándar de Greicodex en el repositorio.",
		Long: `Configura Git para usar el directorio '.githooks' del repositorio,
asegurando que todos los desarrolladores usen los mismos hooks.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := "."
			if len(args) > 0 {
				targetPath = args[0]
			}

			if err := hooksService.InstallHooks(targetPath); err != nil {
				return fmt.Errorf("error al instalar los Git hooks: %w", err)
			}

			fmt.Println("Git hooks instalados exitosamente.")
			return nil
		},
	}
}
