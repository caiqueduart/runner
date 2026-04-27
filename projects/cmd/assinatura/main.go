package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.4"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "assinatura",
		Short: "Ferramenta CLI para assinatura",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual do CLI",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Assinatura CLI - Versão: %s\n", version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
