package cli

import (
	"fmt"
	"grei-cli/internal/adapters/filesystem"
	"grei-cli/internal/adapters/git"
	"grei-cli/internal/core/initializer"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Inicializa un nuevo proyecto con la estructura estándar de Greicodex.",
	Long: `Crea una estructura de directorios estándar y copia las plantillas
necesarias para un nuevo proyecto, incluyendo README, pipelines, docker-compose,
y configuraciones de Helm e IaC.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fsRepo := filesystem.NewRepository()
		gitRepo := git.NewRepository()
		initService := initializer.NewService(fsRepo, gitRepo)

		gitInit, _ := cmd.Flags().GetBool("git")

		err := initService.InitializeProject(targetPath, gitInit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al inicializar el proyecto: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Proyecto inicializado exitosamente en: %s\n", targetPath)
	},
}

func AddInitCommand(root *cobra.Command) {
	initCmd.Flags().Bool("git", false, "Inicializa un repositorio de Git con ramas main y develop.")
	root.AddCommand(initCmd)
}
