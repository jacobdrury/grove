package git

import (
	"context"
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrWorkTreeNotFound       = errors.New("not found")
	workTreeListFormatPattern = regexp.MustCompile(`^(.*)\s+([a-f0-9]+)\s+\[(.*)\]`)
)

// ExecuteWorkTree runs a `git worktree` command with the specified arguments.
func ExecuteWorkTree(ctx context.Context, fmt string, args ...any) (string, error) {
	return execute(ctx, "worktree %v", args...)
}

func ListWorkTrees(ctx context.Context) ([]WorkTree, error) {
	output, err := ExecuteWorkTree(ctx, "list")
	if err != nil {
		return nil, err
	}

	wts := strings.Split(output, "\n")

	return lo.FilterMap(wts, func(v string, _ int) (WorkTree, bool) {
		wt := &WorkTree{}
		err := wt.Scan(v)

		return *wt, err == nil
	}), nil
}

func FindWorkTree(ctx context.Context, branch string) (*WorkTree, error) {
	wts, err := ListWorkTrees(ctx)
	if err != nil {
		return nil, err
	}

	if wt, ok := lo.Find(wts, func(wt WorkTree) bool {
		return wt.Branch == branch
	}); ok {
		return &wt, nil
	}

	return nil, ErrWorkTreeNotFound
}

func CreateWorkTreeFromBranch(ctx context.Context, worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := ExecuteWorkTree(ctx, "add %v %v", worktreePath, branch)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(ctx, branch)
}

func CreateWorkTreeFromNewBranch(ctx context.Context, worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := ExecuteWorkTree(ctx, "add -b %v %v main", branch, worktreePath)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(ctx, branch)
}
