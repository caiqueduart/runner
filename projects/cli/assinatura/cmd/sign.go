package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Assinando documento... wip...")
	},
}

func init() {
	RootCmd.AddCommand(signCmd)
}
