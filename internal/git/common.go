package git

import (
	"context"
	"strings"

	"github.com/samber/lo"
)

func Pull(ctx context.Context) error {
	_, err := execute(ctx, "pull")
	return err
}

func Fetch(ctx context.Context, args ...string) error {
	_, err := execute(ctx, "fetch %v", strings.Join(args, " "))
	return err
}

func BranchExists(ctx context.Context, name string) bool {
	output, err := execute(ctx, "branch --list %v", name)
	if err != nil {
		return false
	}

	if len(strings.TrimSpace(output)) > 0 {
		return true
	}

	output, err = execute(ctx, "ls-remote --heads origin %v", name)
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(output)) > 0
}

func ListBranches(ctx context.Context) ([]string, error) {
	output, err := execute(ctx, "for-each-ref --format='%%(refname:short)' refs/heads/ refs/remotes/")
	if err != nil {
		return nil, err
	}

	branches := strings.Split(output, "\n")
	branches = lo.Map(branches, func(b string, _ int) string {
		b = strings.Trim(b, " '")
		b = strings.TrimPrefix(b, "origin/")

		return b
	})
	branches = lo.Uniq(branches)
	branches = lo.Filter(branches, func(b string, _ int) bool {
		return len(b) > 0 && b != "remote"
	})

	return branches, nil
}
