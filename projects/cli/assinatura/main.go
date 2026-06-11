package main

import (
	"bufio"
	"os"
	"path/filepath"
	"runner/assinatura/cmd"
	"strings"
)

func main() {
	loadNativeEnv()

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadNativeEnv() {
	// Procura pelo .env subindo até 5 níveis (para suportar execução de subpastas)
	currDir, _ := os.Getwd()
	envName := ".env"

	for i := 0; i < 5; i++ {
		envPath := filepath.Join(currDir, envName)
		file, err := os.Open(envPath)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
			return // Sucesso
		}

		// Sobe um nível
		parentDir := filepath.Dir(currDir)
		if parentDir == currDir {
			break // Chegou na raiz do sistema
		}
		currDir = parentDir
	}
}
