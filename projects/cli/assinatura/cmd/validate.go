package cmd

import (
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Valida a assinatura de um documento",
	Args:  cobra.MaximumNArgs(1),
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		runValidate()
	},
}

func runValidate() {
	// os.Args[2:] pega tudo após 'validate'
	args := os.Args[2:]

	output, err := internal.ExecJavaSigner("validate", args)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(output)
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
