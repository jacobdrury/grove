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

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	defer func() {
		err = os.Chdir(wd)
		if err != nil {
			slog.Error("error restoring working directory", slog.String("workingDirectory", wd), slog.String("error", err.Error()))
		}
	}()

	return f()
}
