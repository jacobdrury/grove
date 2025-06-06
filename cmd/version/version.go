package version

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	ErrBuildInfoUnavailable = errors.New("build information unavailable")
)

var Command = &cobra.Command{
	Use:     "version",
	Aliases: []string{"-v", "--version"},
	Short:   "Gets the current version of Grove",
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ErrBuildInfoUnavailable
	}

	fmt.Println(info.Main.Version)
	return nil
}
