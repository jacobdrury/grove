package util

import "os"

func InDirectory(dir string, f func() error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	os.Chdir(dir)
	defer os.Chdir(wd)
	return f()
}
