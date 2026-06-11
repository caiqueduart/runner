package cmd

import (
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	sFile string
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		runSign()
	},
}

func runSign() {
	// Pega todos os argumentos após 'sign'
	// os.Args[0] = cli, os.Args[1] = sign, os.Args[2:] = restante
	args := os.Args[2:]

	output, err := internal.ExecJavaSigner("sign", args)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(output)
}

func init() {
	signCmd.Flags().StringVar(&sFile, "file", "", "Caminho do arquivo para assinatura")

	RootCmd.AddCommand(signCmd)
}
