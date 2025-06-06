package git

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrWorkTreeNotFound = errors.New("not found")
)

var (
	WorkTreeListFormatPattern = regexp.MustCompile(`^(.*)\s+([a-f0-9]+)\s+\[(.*)\]`)
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
		matches := WorkTreeListFormatPattern.FindStringSubmatch(val)
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

func ValidateGitInstallation() error {
	_, err := exec.LookPath("git")
	return err
}

func IsGitRepository(ctx context.Context) (bool, error) {
	res, err := execute(ctx, "rev-parse --is-inside-work-tree")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(res) == "true", nil
}

func ListWorkTrees(ctx context.Context) ([]WorkTree, error) {
	output, err := execute(ctx, "worktree list")
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

func Pull(ctx context.Context) error {
	_, err := execute(ctx, "pull")
	return err
}

func Fetch(ctx context.Context, args ...string) error {
	_, err := execute(ctx, "fetch %v", strings.Join(args, " "))
	return err
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

// git worktree add $newWorkTree $branch_name
func CreateWorkTreeFromBranch(ctx context.Context, worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := execute(ctx, "worktree add %v %v", worktreePath, branch)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(ctx, branch)
}

func CreateWorkTreeFromNewBranch(ctx context.Context, worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := execute(ctx, "worktree add -b %v %v main", branch, worktreePath)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(ctx, branch)
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

func ExecuteWorkTree(ctx context.Context, args string) (string, error) {
	return execute(ctx, "worktree %v", args)
}

func execute(ctx context.Context, format string, args ...any) (string, error) {
	slog.Debug("executing git command", slog.String("command", fmt.Sprintf(format, args...)))

	cmdFormatted := fmt.Sprintf(format, args...)
	cmd := exec.CommandContext(ctx, "git", strings.Split(cmdFormatted, " ")...)
	output, err := cmd.CombinedOutput()

	slog.Debug("git command output", slog.String("command", cmdFormatted), slog.String("output", string(output)))

	return string(output), err
}
