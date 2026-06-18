package cmd

import (
	"fmt"
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

// startCmd representa o comando para iniciar o simulador.
// Ele executa o provisionamento automático (download) caso o simulador ou o JDK não existam.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o simulador HubSaúde em background",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.EnsureSimuladorRunning(); err != nil {
			return fmt.Errorf("erro ao iniciar simulador: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}
