package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var file string

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",

	Run: func(cmd *cobra.Command, args []string) {

		if file == "" {
			fmt.Println("Erro: Use a flag -f para indicar o caminho do arquivo.")
			return
		}

		// 1. Validação do CLI: O arquivo do usuário existe?
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Erro ao ler o arquivo '%s': %v\n", file, err)
			return
		}

		const JAR_PATH = "../../assinador/src/Main.java"
		var fullArgs = []string{ /* "-jar", */ JAR_PATH, string(content)}

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
	signCmd.Flags().StringVarP(&file, "file", "f", "", "Arquivo para assinatura")

	RootCmd.AddCommand(signCmd)
}
