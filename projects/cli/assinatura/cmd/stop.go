package cmd

import (
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o servidor do assinador",
	Run: func(cmd *cobra.Command, args []string) {
		runStop()
	},
}

func runStop() {
	internal.LogFeedback("ASSINATURA SERVIDOR", "Encerrando servidor...")
	_, err := internal.CallJavaServer("stop", "")
	if err != nil {
		internal.LogFeedback("ASSINATURA SERVIDOR", "Aviso: Servidor já está desligado ou inacessível.")
		internal.ClearPIDFile()
		return
	}
	internal.ClearPIDFile()
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
