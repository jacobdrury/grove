package util

import (
	"log/slog"
	"os"
)

func InDirectory[T any](dir string, f func() (T, error)) (T, error) {
	var zero T

	wd, err := os.Getwd()
	if err != nil {
		return zero, err
	}

	slog.Debug("changing directory", slog.String("path", dir))
	err = os.Chdir(dir)
	if err != nil {
		return zero, err
	}

	defer func() {
		slog.Debug("changing directory", slog.String("path", wd))
		err = os.Chdir(wd)
		if err != nil {
			slog.Error("error restoring working directory", slog.String("workingDirectory", wd), slog.String("error", err.Error()))
		}
	}()

	result, err := f()
	if err != nil {
		return zero, err
	}

	return result, nil
}

func InDirectoryNoResult(dir string, f func() error) error {
	_, err := InDirectory(dir, func() (any, error) {
		return nil, f()
	})

	return err
}
