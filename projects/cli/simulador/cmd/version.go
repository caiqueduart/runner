package cmd

import (
	"fmt"
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão da CLI do Simulador",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Simulador CLI v%s\n", internal.SimuladorCLIVersion)
		fmt.Printf("Compatível com Simulador JAR v%s\n", internal.CompatibleSimuladorVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
