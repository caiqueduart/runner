package cmd

import (
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o simulador HubSaúde em background",
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.EnsureSimuladorRunning(); err != nil {
			internal.PrintError("Erro ao iniciar simulador: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
