package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jacobdrury/grove/cmd/checkout"
	"github.com/jacobdrury/grove/cmd/initialize"
	"github.com/jacobdrury/grove/cmd/version"
	"github.com/jacobdrury/grove/internal/git"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "wt",
	Short:         "wt is a git worktree manager",
	SilenceErrors: true, // Errors are output to stderr so we don't need to print them
	Args: func(cmd *cobra.Command, args []string) error {
		// If there are args, check if the first is a valid subcommand
		if len(args) < 1 {
			return nil
		}

		// If the first arg is a subcommand, return
		if lo.SomeBy(cmd.Commands(), func(c *cobra.Command) bool {
			return c.Name() == args[0] || c.HasAlias(args[0])
		}) {
			return nil
		}

		slog.Debug("no subcommand found, passing through args to git worktree command", slog.String("args", strings.Join(args, " ")))

		output, err := git.ExecuteWorkTree(cmd.Context(), strings.Join(args, " "))
		if err != nil {
			return err
		}

		fmt.Print(output)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}

		return nil
	},
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		func() {
			err := git.ValidateGitInstallation()
			cobra.CheckErr(err)
		},
	)

	rootCmd.AddCommand(
		checkout.Command,
		initialize.Command,
		version.Command,
	)

	rootCmd.SilenceUsage = true
}
