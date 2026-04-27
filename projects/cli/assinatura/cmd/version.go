package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CLI_VERSION = "0.1.6"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Assinatura CLI - Versão: %s\n", CLI_VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
