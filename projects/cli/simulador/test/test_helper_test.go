package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

type commandResult struct {
	stdout   string
	stderr   string
	exitCode int
}

var simulatorBinary string

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "simulador-cli-test-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	binaryName := "simulador"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	simulatorBinary = filepath.Join(tempDir, binaryName)

	build := exec.Command("go", "build", "-o", simulatorBinary, "..")
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "falha ao compilar CLI: %v\n%s", err, output)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func runCLI(t *testing.T, baseURL string, args ...string) commandResult {
	t.Helper()

	home := t.TempDir()
	command := exec.Command(simulatorBinary, args...)
	command.Env = append(os.Environ(),
		"SIMULADOR_BASE_URL="+baseURL,
		"HOME="+home,
		"USERPROFILE="+home,
	)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	result := commandResult{stdout: stdout.String(), stderr: stderr.String()}
	if exitError, ok := err.(*exec.ExitError); ok {
		result.exitCode = exitError.ExitCode()
	} else if err != nil {
		t.Fatalf("falha ao executar CLI: %v", err)
	}

	return result
}
