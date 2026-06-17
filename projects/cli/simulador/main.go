package main

import (
	"bufio"
	"os"
	"path/filepath"
	"runner/simulador/cmd"
	"strings"
)

func main() {
	loadNativeEnv()

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// tenta carregar um arquivo .env subindo até 5 níveis de diretório. (modo desenvolvedor)
func loadNativeEnv() {
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
			return // sucesso
		}

		parentDir := filepath.Dir(currDir)

		if parentDir == currDir {
			break // raiz do sistema
		}

		currDir = parentDir
	}
}
