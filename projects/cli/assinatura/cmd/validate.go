package cmd

import (
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	vFile   string
	vCmdKey = "validate"
)

var validateCmd = &cobra.Command{
	Use:   vCmdKey,
	Short: "Valida um documento",
	Run: func(cmd *cobra.Command, args []string) {
		runValidate()
	},
}

func runValidate() {

	// Leitura das informações do arquivo passado nos parâmetros
	fileInfo, err := os.Stat(vFile)
	if err != nil {
		internal.PrintError("Erro ao ler o arquivo '%s': \n%v", vFile, err)
		return
	}

	// Execução do comando para o assinador
	output, err := internal.ExecJavaSigner(string(fileInfo.Name()), vCmdKey)
	if err != nil {
		internal.PrintError("Erro ao executar o assinador: \n%v", err)
		return
	}

	fmt.Println(output)
}

func init() {
	validateCmd.Flags().StringVarP(&vFile, "file", "f", "", "Arquivo para assinatura")
	validateCmd.MarkFlagRequired("file")

	RootCmd.AddCommand(validateCmd)
}
