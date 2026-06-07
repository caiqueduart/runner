package cmd

import (
	"fmt"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Assinatura CLI - Versão: %s\n", internal.AssinaturaCLIVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
