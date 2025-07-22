package checkout

import (
	"fmt"

	"github.com/jacobdrury/grove/internal/config"
	"github.com/jacobdrury/grove/internal/grove"
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

var (
	pipe    bool
	noHooks bool
)

func init() {
	Command.Flags().BoolVarP(&pipe, "pipe", "p", false, "pipe worktree path to stdout")
	Command.Flags().BoolVarP(&noHooks, "no-hooks", "n", false, "do not run hooks")
}

func run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if noHooks {
		ctx = config.ContextWithNoHooks(ctx)
	}

	if pipe {
		ctx = config.ContextWithPipe(ctx)
	}

	g, err := grove.GetInstance()
	if err != nil {
		return err
	}

	wt, err := g.Checkout(ctx, grove.CheckoutArgs{
		Branch: args[0],
	})
	if err != nil {
		return err
	}

	if config.Pipe(ctx) {
		_, err := fmt.Fprint(cmd.OutOrStdout(), wt.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	return grove.Load()
}
