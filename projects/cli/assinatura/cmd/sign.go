package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var file string

const (
	CMD_KEY  = "sign"
	JAR_PATH = "../../assinador/src/Main.java"
)

var signCmd = &cobra.Command{
	Use:   CMD_KEY,
	Short: "Assina um documento",
	Run: func(cmd *cobra.Command, args []string) {
		runSign()
	},
}

func runSign() {
	// Validação de existência do assinador
	if _, err := os.Stat(JAR_PATH); os.IsNotExist(err) {
		internal.PrintError("Erro: O Assinador não foi encontrado em '%s'", JAR_PATH)
		return
	}

	// Leitura do arquivo passado nos parâmetros
	content, err := os.ReadFile(file)
	if err != nil {
		internal.PrintError("Erro ao ler o arquivo '%s': \n%v", file, err)
		return
	}

	// Execução do comando externo
	output, err := execJavaSigner(string(content))
	if err != nil {
		internal.PrintError("Erro ao executar o assinador: \n%v", err)
		return
	}

	fmt.Println(output)
}

func execJavaSigner(payload string) (string, error) {
	javaCmd := exec.Command("java", JAR_PATH, CMD_KEY, payload)

	output, err := javaCmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(output), nil
}

func init() {
	signCmd.Flags().StringVarP(&file, "file", "f", "", "Arquivo para assinatura")
	signCmd.MarkFlagRequired("file")

	RootCmd.AddCommand(signCmd)
}
