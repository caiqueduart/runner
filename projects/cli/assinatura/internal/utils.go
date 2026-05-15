package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const jarPath = "../../assinador/src/Main.java"

func PrintError(format string, a ...any) {
	fmt.Printf(ColorRed+format+ColorReset, a...)
}

// GetJavaPath retorna o caminho para o executável java.
// Primeiro verifica se está no PATH, depois no diretório gerenciado ~/.hubsaude/jdk.
func GetJavaPath() (string, error) {
	var foundPath string
	var foundVersion bool

	// 1. Verifica no PATH
	path, err := exec.LookPath("java")
	if err == nil {
		if isJava21(path) {
			return path, nil
		}
		foundPath = path
		foundVersion = true
	}

	// 2. Verifica no diretório gerenciado ~/.hubsaude/jdk
	home, err := os.UserHomeDir()
	if err == nil {
		javaBin := "java"
		if runtime.GOOS == "windows" {
			javaBin = "java.exe"
		}

		managedPath := filepath.Join(home, ".hubsaude", "jdk", "bin", javaBin)
		if _, err := os.Stat(managedPath); err == nil {
			if isJava21(managedPath) {
				return managedPath, nil
			}
			foundPath = managedPath
			foundVersion = true
		}
	}

	if foundVersion {
		return "", fmt.Errorf("Java encontrado em '%s', mas não é a versão 21 exigida", foundPath)
	}

	return "", fmt.Errorf("JDK 21 não encontrado no PATH ou em ~/.hubsaude/jdk")
}

// isJava21 verifica se o binário java informado é da versão 21.
func isJava21(javaPath string) bool {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	// A saída do java -version costuma ir para stderr e contém algo como "openjdk version \"21.0.x\""
	versionOutput := string(output)
	return strings.Contains(versionOutput, "version \"21")
}

func ExecJavaSigner(fileName string, cmdKey string) (string, error) {

	// Validação de existência do assinador
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return "", fmt.Errorf("O Assinador não foi encontrado em '%s'.", jarPath)
	}

	javaPath, err := GetJavaPath()
	if err != nil {
		return "", err
	}

	javaCmd := exec.Command(javaPath, jarPath, cmdKey, fileName)

	output, err := javaCmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(output), nil
}
