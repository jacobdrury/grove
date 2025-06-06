package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jacobdrury/wt/cmd"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

func init() {
	var writer io.Writer = os.Stdout
	if runtime.GOOS == "windows" {
		writer = colorable.NewColorableStdout()
	}

	logLevel := slog.LevelInfo
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		err := logLevel.UnmarshalText([]byte(level))
		if err != nil {
			panic(fmt.Sprintf("invalid log level %s: %v", level, err))
		}
	}

	slog.SetDefault(slog.New(tint.NewHandler(writer, &tint.Options{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the time attribute to de-clutter the output
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}

			return a
		},
	})))
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer stop()

	cmd.Execute(ctx)
}
