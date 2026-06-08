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
	envPath := filepath.Join("..", "..", "..", ".env")
	file, err := os.Open(envPath)

	if err != nil {
		return
	}

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
}
