package cmd

import (
	"fmt"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Valida a assinatura de um documento",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vFile := ""
		if len(args) > 0 {
			vFile = args[0]
		}
		runValidate(vFile)
	},
}

func runValidate(file string) {
	output, err := internal.ExecJavaSigner(file, "validate")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(output)
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
