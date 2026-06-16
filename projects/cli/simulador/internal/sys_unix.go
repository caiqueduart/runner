//go:build !windows

package internal

import (
	"os/exec"
)

// SetDetachedProcess configura o comando para rodar de forma independente em sistemas Unix.
func SetDetachedProcess(cmd *exec.Cmd) {
	// Em sistemas Unix, o comportamento padrão do cmd.Start() com Release()
	// costuma ser suficiente para processos em background.
}
