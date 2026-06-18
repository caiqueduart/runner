package cmd

import (
	"fmt"
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica o status do simulador HubSaúde",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := internal.GetSimuladorStatus()
		if err != nil {
			internal.LogFeedback("SIMULADOR CONFIG", "Simulador não está respondendo corretamente (Offline ou erro).")
			return fmt.Errorf("falha ao consultar simulador: %w", err)
		}

		internal.LogFeedback("SIMULADOR SERVIDOR", "Simulador Online!")
		fmt.Printf("Detalhes: %s\n", info)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
