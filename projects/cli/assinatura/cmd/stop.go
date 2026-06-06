package cmd

import (
	"fmt"
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
	fmt.Println("Encerrando o servidor do assinador...")
	output, err := internal.CallJavaServer("stop", "")
	if err != nil {
		fmt.Printf("Aviso: Não foi possível encerrar o servidor (provavelmente já está desligado): %v\n", err)
		return
	}
	fmt.Println(output)
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
