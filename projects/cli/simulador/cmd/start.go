package cmd

import (
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

// startCmd representa o comando para iniciar o simulador.
// Ele executa o provisionamento automático (download) caso o simulador ou o JDK não existam.
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
