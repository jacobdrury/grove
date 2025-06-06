package git

import (
	"fmt"
	"strings"
)

type WorkTree struct {
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
		if len(matches) != 4 {
			return fmt.Errorf("invalid worktree format")
		}

		w.Path = strings.TrimSpace(matches[1])
		w.Head = strings.TrimSpace(matches[2])
		w.Branch = strings.TrimSpace(matches[3])
	default:
		return fmt.Errorf("invalid scan type")
	}

	return nil
}
