package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CmdResult armazena o resultado da execução de um comando CLI
type CmdResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// RunCLI executa a CLI de assinatura com os argumentos fornecidos
func RunCLI(args ...string) (*CmdResult, error) {
	// Caminho para o main.go (estamos em projects/cli/assinatura/test)
	mainPath := filepath.Join("..", "main.go")

	// Prepara o comando: go run ../main.go [args]
	cmdArgs := append([]string{"run", mainPath}, args...)
	cmd := exec.Command("go", cmdArgs...)

	// Captura buffers
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Executa
	err := cmd.Run()

	result := &CmdResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			return nil, err
		}
	} else {
		result.ExitCode = 0
	}

	return result, nil
}

// CreateDummyFile cria um arquivo temporário para testes e retorna seu caminho absoluto
func CreateDummyFile(name, content string) (string, error) {
	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		return "", err
	}
	absPath, _ := filepath.Abs(name)
	return absPath, nil
}

// CleanUpFiles remove arquivos gerados pelos testes
func CleanUpFiles(patterns ...string) {
	for _, pattern := range patterns {
		files, _ := filepath.Glob(pattern)
		for _, f := range files {
			os.Remove(f)
		}
	}
}

// Contains checks if a string contains another, ignoring ANSI escape codes
func ContainsIgnoringColors(str, substr string) bool {
	cleanStr := stripColors(str)
	return strings.Contains(cleanStr, substr)
}

func stripColors(str string) string {
	// Simplificado: remove prefixos comuns de cores ANSI
	replacer := strings.NewReplacer(
		"\033[31m", "",
		"\033[32m", "",
		"\033[33m", "",
		"\033[34m", "",
		"\033[36m", "",
		"\033[0m", "",
		"\u001B[31m", "",
		"\u001B[32m", "",
		"\u001B[33m", "",
		"\u001B[34m", "",
		"\u001B[36m", "",
		"\u001B[0m", "",
	)
	return replacer.Replace(str)
}
