package cmd

import (
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o simulador HubSaúde",
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.StopSimulador(); err != nil {
			internal.PrintError("Erro ao parar simulador: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
