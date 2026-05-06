package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var file string

const CMD_KEY = "sign"

var signCmd = &cobra.Command{
	Use:   CMD_KEY,
	Short: "Assina um documento",

	Run: func(cmd *cobra.Command, args []string) {

		content, err := os.ReadFile(file)
		if err != nil {
			internal.PrintError("Ocorreu um erro ao ler o arquivo '%s': \n'%s'", file, err)
			return
		}

		const JAR_PATH = "../../assinador/src/Main.java"
		var fullArgs = []string{ /* "-jar", */ JAR_PATH, CMD_KEY, string(content)}

		if _, err := os.Stat(JAR_PATH); os.IsNotExist(err) {
			internal.PrintError("Erro: O arquivo '%s' não foi encontrado.", JAR_PATH)
			return
		}

		javaCmd := exec.Command("java", fullArgs...)

		output, err := javaCmd.CombinedOutput()

		if err != nil {
			internal.PrintError("Ocorreu um erro ao executar o assinador: \n'%s'", err)
			return
		}

		fmt.Println(string(output))
	},
}

func init() {
	signCmd.Flags().StringVarP(&file, "file", "f", "", "Arquivo para assinatura")

	RootCmd.AddCommand(signCmd)
}
