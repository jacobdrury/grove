package wt

import (
	"errors"
	"fmt"
	"os"

	"github.com/jacobdrury/wt/internal/git"
	"github.com/jacobdrury/wt/internal/util"
)

type CheckoutArgs struct {
	Branch string // Supports aliases j/fm-3311
}

func Checkout(arg CheckoutArgs) error {
	fmt.Println("Checking out", arg.Branch)

	ctx, err := GetContext()
	if err != nil {
		return err
	}

	branch := ctx.ResolveBranch(arg.Branch)

	wt, err := git.FindWorkTree(branch)
	if err != nil && !errors.Is(err, git.ErrWorkTreeNotFound) {
		return err
	}

	if wt != nil {
		fmt.Println("Worktree already exists, updating it")
		return checkoutWorkTree(ctx, wt)
	}

	err = git.Fetch("-p")
	if err != nil {
		return err
	}

	// 2. If branch exists on remote, add a new worktree for
	if git.BranchExists(branch) {
		fmt.Println("Branch exists on remote, creating worktree from remote")

		wt, err = git.CreateWorkTreeFromBranch(ctx.Config.WorkTreesDirectory, branch)
		if err != nil {
			return err
		}

		return checkoutWorkTree(ctx, wt)
	}

	mainWt, err := git.FindWorkTree("main")
	if err != nil {
		return fmt.Errorf("error finding main worktree: %v", err)
	}

	// Update main worktree
	err = util.InDirectory(mainWt.Path, func() error {
		fmt.Println("Updating main worktree")

		return git.Pull()
	})
	if err != nil {
		return err
	}

	fmt.Println("Creating new worktree from new branch")
	wt, err = git.CreateWorkTreeFromNewBranch(ctx.Config.WorkTreesDirectory, branch)
	if err != nil {
		return err
	}

	return checkoutWorkTree(ctx, wt)
}

func checkoutWorkTree(ctx *WorkTreeContext, wt *git.WorkTree) error {
	err := os.Chdir(wt.Path)
	if err != nil {
		return err
	}

	// We don't care if it fails, just want to try and update the branch
	_ = git.Pull()

	// Copy seed files
	err = ctx.SeedWorkTree(wt)
	if err != nil {
		return err
	}

	return ctx.ExecuteAfterCheckoutHooks()
}
