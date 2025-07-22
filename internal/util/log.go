package util

import (
	"context"
	"log/slog"

	"github.com/jacobdrury/grove/internal/config"
)

func LogInfo(ctx context.Context, msg string, args ...any) {
	if config.Pipe(ctx) {
		// No-op
		return
	}

	slog.InfoContext(ctx, msg, args...)
}
