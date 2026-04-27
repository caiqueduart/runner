package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida um documento",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Validando um documento... wip")
	},
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
