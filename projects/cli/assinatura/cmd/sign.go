package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",

	Run: func(cmd *cobra.Command, args []string) {
		// Verifica e o arquivo .jar está
		const JAR_PATH = "../../assinador/dist/assinador.jar"

		if _, err := os.Stat(JAR_PATH); os.IsNotExist(err) {
			fmt.Printf("Erro: O arquivo %s não foi encontrado.", JAR_PATH)
			return
		}

		javaCmd := exec.Command("java", "-jar", JAR_PATH)

		output, err := javaCmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Erro ao executar o assinador: %s\n", err)
			return
		}

		fmt.Println("Resultado do Assinador:")
		fmt.Println(string(output))
	},
}

func init() {
	RootCmd.AddCommand(signCmd)
}
