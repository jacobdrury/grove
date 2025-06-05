package git

import (
	"errors"
	"fmt"
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

func IsGitRepository() (bool, error) {
	res, err := execute("rev-parse --is-inside-work-tree")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(res) == "true", nil
}

func ListWorkTrees() ([]WorkTree, error) {
	output, err := execute("worktree list")
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

func Pull() error {
	_, err := execute("pull")
	return err
}

func Fetch(args ...string) error {
	_, err := execute("fetch %v", strings.Join(args, " "))
	return err
}

func FindWorkTree(branch string) (*WorkTree, error) {
	wts, err := ListWorkTrees()
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
func CreateWorkTreeFromBranch(worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := execute("worktree add %v %v", worktreePath, branch)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(branch)
}

func CreateWorkTreeFromNewBranch(worktreesPath string, branch string) (*WorkTree, error) {
	worktreePath := filepath.Join(worktreesPath, branch)

	_, err := execute("worktree add -b %v %v main", branch, worktreePath)
	if err != nil {
		return nil, err
	}

	return FindWorkTree(branch)
}

func BranchExists(name string) bool {
	output, err := execute("branch --list %v", name)
	if err != nil {
		return false
	}

	if len(strings.TrimSpace(output)) > 0 {
		return true
	}

	output, err = execute("ls-remote --heads origin %v", name)
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(output)) > 0
}

func ExecuteWorkTree(args string) (string, error) {
	return execute("worktree %v", args)
}

func execute(format string, args ...any) (string, error) {
	cmdFormatted := fmt.Sprintf(format, args...)
	cmd := exec.Command("git", strings.Split(cmdFormatted, " ")...)
	output, err := cmd.CombinedOutput()

	// fmt.Printf("%v: %v\n", cmdFormatted, err)
	return string(output), err
}
