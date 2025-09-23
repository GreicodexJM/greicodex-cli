package cli

import (
	"fmt"
	"grei-cli/internal/adapters/syschecker"
	"grei-cli/internal/core/doctor"
	"os"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verifica que el entorno de desarrollo local esté configurado correctamente.",
	Long: `Comprueba que todas las herramientas necesarias (git, docker, helm, etc.)
estén instaladas y disponibles en el PATH del sistema.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🩺 Verificando el entorno de desarrollo...")

		sysChecker := syschecker.New()
		doctorService := doctor.NewService(sysChecker)
		results := doctorService.CheckEnvironment()

		allGood := true
		for _, result := range results {
			if result.Found {
				fmt.Printf("  [✓] %s\n", result.Command)
			} else {
				fmt.Printf("  [✗] %s No encontrado en su sistema. Debe instalarlo.\n", result.Command)
				allGood = false
			}
		}

		if !allGood {
			fmt.Println("\nAlgunas herramientas requeridas no fueron encontradas. Por favor, instálalas y asegúrate de que estén en tu PATH.")
			os.Exit(1)
		}

		fmt.Println("\n¡Tu entorno está listo!")
	},
}

func AddDoctorCommand(root *cobra.Command) {
	root.AddCommand(doctorCmd)
}
