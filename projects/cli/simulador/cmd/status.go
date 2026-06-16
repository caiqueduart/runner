package cmd

import (
	"fmt"
	"runner/simulador/internal"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica o status do simulador HubSaúde",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := internal.GetSimuladorStatus()
		if err != nil {
			internal.LogFeedback("SIMULADOR CONFIG", "Simulador não está respondendo corretamente (Offline ou erro).")
			return
		}

		internal.LogFeedback("SIMULADOR SERVIDOR", "Simulador Online!")
		fmt.Printf("Detalhes: %s\n", info)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
