package internal

import (
	"fmt"
	"os"
	"os/exec"
)

const jarPath = "../../assinador/src/Main.java"

func PrintError(format string, a ...any) {
	fmt.Printf(ColorRed+format+ColorReset, a...)
}

func ExecJavaSigner(fileName string, cmdKey string) (string, error) {

	// Validação de existência do assinador
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return "", fmt.Errorf("O Assinador não foi encontrado em '%s'.", jarPath)
	}

	javaCmd := exec.Command("java", jarPath, cmdKey, fileName)

	output, err := javaCmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(output), nil
}
