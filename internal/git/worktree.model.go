package git

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	workTreeListFormatPattern = regexp.MustCompile(`^((.*)(\/|\\)([^\/\\]+))\s+([a-f0-9]+)\s+\[(.*)\]$`)
)

type WorkTree struct {
	Name   string
	Path   string
	Head   string
	Branch string
}

func (w WorkTree) String() string {
	return fmt.Sprintf("%v %v [%v]", w.Path, w.Head, w.Branch)
}

func (w *WorkTree) Scan(v any) error {
	switch val := v.(type) {
	case string:
		matches := workTreeListFormatPattern.FindStringSubmatch(val)
		if len(matches) != 7 {
			return fmt.Errorf("invalid worktree format")
		}

		w.Path = strings.TrimSpace(matches[1])
		w.Name = strings.TrimSpace(matches[4])
		w.Head = strings.TrimSpace(matches[5])
		w.Branch = strings.TrimSpace(matches[6])
	default:
		return fmt.Errorf("invalid scan type")
	}

	return nil
}
