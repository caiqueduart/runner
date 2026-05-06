package cmd

import (
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	sFile   string
	sCmdKey = "sign"
)

var signCmd = &cobra.Command{
	Use:   sCmdKey,
	Short: "Assina um documento",
	Run: func(cmd *cobra.Command, args []string) {
		runSign()
	},
}

func runSign() {

	// Leitura das informações do arquivo passado nos parâmetros
	fileInfo, err := os.Stat(sFile)
	if err != nil {
		internal.PrintError("Erro ao ler o arquivo '%s': \n%v", sFile, err)
		return
	}

	// Execução do comando para o assinador
	output, err := internal.ExecJavaSigner(string(fileInfo.Name()), sCmdKey)
	if err != nil {
		internal.PrintError("Erro ao executar o assinador: \n%v", err)
		return
	}

	fmt.Println(output)
}

func init() {
	signCmd.Flags().StringVarP(&sFile, "file", "f", "", "Arquivo para assinatura")
	signCmd.MarkFlagRequired("file")

	RootCmd.AddCommand(signCmd)
}
