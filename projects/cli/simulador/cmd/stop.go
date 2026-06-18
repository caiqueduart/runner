package cmd

import (
	"fmt"
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

// stopCmd representa o comando para encerrar o simulador graciosamente.
// Ele utiliza o endpoint /shutdown via HTTPS conforme especificado na US-03.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o simulador HubSaúde",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.StopSimulador(); err != nil {
			return fmt.Errorf("erro ao parar simulador: %w", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
