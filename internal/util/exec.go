package util

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExecShellCmd(ctx context.Context, shell string, cmd string) (string, error) {
	var command *exec.Cmd

	// Normalize shell name for comparison
	shellBase := strings.ToLower(filepath.Base(shell))

	switch shellBase {
	case "powershell", "pwsh":
		// PowerShell: use -Command
		command = exec.CommandContext(ctx, shell, "-Command", cmd)
	case "cmd", "cmd.exe":
		// cmd.exe: use /C
		command = exec.CommandContext(ctx, shell, "/C", cmd)
	default:
		// Unix shells: use -c
		command = exec.CommandContext(ctx, shell, "-c", cmd)
	}

	out, err := command.CombinedOutput()
	return string(out), err
}
