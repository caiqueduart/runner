package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runner/assinatura/internal"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Valida a assinatura de um documento",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runValidate(args)
	},
}

func runValidate(args []string) {
	fileName := ""
	if len(args) == 1 {
		fileName = args[0]
	}

	output, err := internal.ExecJavaSigner("validate", fileName, executionOptions())

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
	var res internal.ValidationResponse
	if err := json.Unmarshal([]byte(output), &res); err == nil && res.Code != "" {
		fmt.Printf("\x1b[36m[ASSINATURA]\x1b[0m %s\n", res.Message)
		fmt.Printf("\x1b[36m[ASSINATURA]\x1b[0m O Arquivo '\x1b[32m%s\x1b[0m' está assinado sob o código '\x1b[32m%s\x1b[0m'.\n", res.FileName, res.Code)
		return
	}

	// se não for JSON ou falhar, imprime o output bruto
	fmt.Print(output)
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
