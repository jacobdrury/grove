package util

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func ExecShellCmd(shell string, cmd string) (string, error) {
	var command *exec.Cmd

	// Normalize shell name for comparison
	shellBase := strings.ToLower(filepath.Base(shell))

	switch shellBase {
	case "powershell", "pwsh":
		// PowerShell: use -Command
		command = exec.Command(shell, "-Command", cmd)
	case "cmd", "cmd.exe":
		// cmd.exe: use /C
		command = exec.Command(shell, "/C", cmd)
	default:
		// Unix shells: use -c
		command = exec.Command(shell, "-c", cmd)
	}

	out, err := command.CombinedOutput()
	return string(out), err
}
