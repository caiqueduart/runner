package cmd

import (
	"fmt"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	sFile string
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",
	Run: func(cmd *cobra.Command, args []string) {
		runSign()
	},
}

func runSign() {
	output, err := internal.ExecJavaSigner(sFile, "sign")
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
