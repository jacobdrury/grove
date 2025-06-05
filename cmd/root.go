package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jacobdrury/wt/cmd/checkout"
	"github.com/jacobdrury/wt/cmd/initialize"
	"github.com/jacobdrury/wt/internal/git"
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

		if lo.SomeBy(cmd.Commands(), func(c *cobra.Command) bool {
			return c.Name() == args[0] || c.HasAlias(args[0])
		}) {
			return nil
		}

		output, err := git.ExecuteWorkTree(strings.Join(args, " "))
		if err != nil {
			return err
		}

		print(output)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	)

	rootCmd.SilenceUsage = true
}
