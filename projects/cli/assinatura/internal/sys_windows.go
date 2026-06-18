//go:build windows

package internal

import (
	"os/exec"
	"syscall"
)

// configura o comando para rodar de forma independente no Windows.
func SetDetachedProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000008 | // DETACHED_PROCESS
			0x00000200, // CREATE_NEW_PROCESS_GROUP
		HideWindow: true,
	}
}
