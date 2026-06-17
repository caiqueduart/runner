//go:build !windows

package internal

import (
	"os/exec"
)

// configura o comando para rodar de forma independente em sistemas Unix.
func SetDetachedProcess(cmd *exec.Cmd) {}
