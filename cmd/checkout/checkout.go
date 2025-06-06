package checkout

import (
	"github.com/jacobdrury/grove/internal/wt"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:               "checkout",
	Aliases:           []string{"create", "co"},
	Short:             "Checkout a branch in a new worktree",
	Args:              cobra.ExactArgs(1),
	RunE:              run,
	PersistentPreRunE: persistentPreRun,
}

func run(cmd *cobra.Command, args []string) error {
	return wt.Checkout(cmd.Context(), wt.CheckoutArgs{
		Branch: args[0],
	})
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	return wt.LoadContext()
}
