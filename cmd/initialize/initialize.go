package initialize

import (
	"github.com/jacobdrury/grove/internal/wt"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new worktree manager in the current directory",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	return wt.CreateContext(cmd.Context())
}
