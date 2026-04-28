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

		const JAR_PATH = "../../assinador/dist/assinador.jar"
		var fullArgs = append([]string{"-jar", JAR_PATH}, args...)

		// Verifica se o arquivo .jar existe
		if _, err := os.Stat(JAR_PATH); os.IsNotExist(err) {
			fmt.Printf("Erro: O arquivo %s não foi encontrado.", JAR_PATH)
			return
		}

		javaCmd := exec.Command("java", fullArgs...)

		output, err := javaCmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Erro ao executar o assinador: %s\n", err)
			return
		}

		fmt.Println(string(output))
	},
}

func init() {
	RootCmd.AddCommand(signCmd)
}
