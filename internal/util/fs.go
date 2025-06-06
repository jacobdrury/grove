package util

import (
	"log/slog"
	"os"
)

func InDirectory(dir string, f func() error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	slog.Debug("changing directory", slog.String("path", dir))
	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	defer func() {
		slog.Debug("changing directory", slog.String("path", wd))
		err = os.Chdir(wd)
		if err != nil {
			slog.Error("error restoring working directory", slog.String("workingDirectory", wd), slog.String("error", err.Error()))
		}
	}()

	return f()
}
