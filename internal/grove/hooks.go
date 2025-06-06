package grove

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jacobdrury/grove/internal/util"
)

func (grove *Grove) executeAfterCheckoutHooks(ctx context.Context) error {
	slog.Debug("executing after checkout hooks", slog.Int("numberOfHooks", len(grove.Config.Hooks.AfterCheckout)))
	for _, hook := range grove.Config.Hooks.AfterCheckout {
		slog.Info("executing hook", slog.String("hook", hook))

		err := util.ExecShellCmd(ctx, grove.Config.Hooks.Shell, hook)
		if err != nil {
			return fmt.Errorf("error executing hook %s: %v", hook, err)
		}
	}

	slog.Debug("after checkout hooks executed")

	return nil
}
