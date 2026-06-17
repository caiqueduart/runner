package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var (
	sFile string
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Assina um documento",
	Run: func(cmd *cobra.Command, args []string) {
		runSign()
	},
}

func runSign() {
	output, err := internal.ExecJavaSigner("sign", sFile, executionOptions())
	if err != nil {
		if jErr, ok := err.(*internal.JavaError); ok {
			if jErr.Type == "user" {
				internal.LogFeedback("ASSINATURA", "Erro do usuário: %s", jErr.Msg)
				os.Exit(1)
			} else {
				internal.LogFeedback("ASSINATURA", "Erro do sistema: %s", jErr.Msg)
				os.Exit(2)
			}
		} else if execErr, ok := err.(*internal.ExecError); ok {
			fmt.Print(execErr.Output)
			os.Exit(execErr.Code)
		}
		internal.LogFeedback("ASSINATURA", "Erro: %v", err)
		os.Exit(2)
	}

	// tenta decodificar como JSON (sucesso)
	var res internal.SignatureResponse
	if err := json.Unmarshal([]byte(output), &res); err == nil && res.Code != "" {
		fmt.Printf("\x1b[36m[ASSINATURA]\x1b[0m %s\n", res.Message)
		fmt.Printf("\x1b[36m[ASSINATURA]\x1b[0m O Arquivo '\x1b[32m%s\x1b[0m' gerou o código de assinatura '\x1b[32m%s\x1b[0m'.\n", res.FileName, res.Code)
		fmt.Printf("\x1b[36m[ASSINATURA]\x1b[0m Arquivo gerado em: \x1b[32m%s\x1b[0m\n", res.SignOutputPath)
		return
	}

	// se não for JSON ou falhar, imprime o output bruto (preserva cores do JAR se houver)
	fmt.Print(output)
}

func init() {
	signCmd.Flags().StringVar(&sFile, "file", "", "Caminho do arquivo para assinatura")

	RootCmd.AddCommand(signCmd)
}
