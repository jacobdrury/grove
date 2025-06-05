package wt

import (
	"github.com/jacobdrury/wt/internal/git"
	"github.com/samber/lo"
)

func ListWorkTrees() error {
	wts, err := git.ListWorkTrees()
	if err != nil {
		return err
	}

	lo.ForEach(wts, func(wt git.WorkTree, _ int) {
		println(wt.String())
	})

	return nil
}
