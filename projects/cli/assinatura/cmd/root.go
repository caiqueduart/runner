package cmd

import (
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	localMode            bool
	serverPort           int
	serverTimeoutMinutes int
)

var RootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "Ferramenta CLI para assinatura digital e validação",
}

func executionOptions() internal.ExecutionOptions {
	return internal.ExecutionOptions{
		Local:          localMode,
		Port:           serverPort,
		TimeoutMinutes: serverTimeoutMinutes,
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&localMode, "local", false, "Executa o assinador.jar diretamente, sem servidor HTTP")
	RootCmd.PersistentFlags().IntVar(&serverPort, "port", internal.DefaultServerPort, "Porta do servidor HTTP do assinador")
	RootCmd.PersistentFlags().IntVar(&serverTimeoutMinutes, "timeout", internal.DefaultServerTimeoutMinutes, "Minutos de inatividade antes do servidor encerrar")
}
