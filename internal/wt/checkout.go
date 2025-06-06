package wt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jacobdrury/wt/internal/git"
	"github.com/jacobdrury/wt/internal/util"
)

type CheckoutArgs struct {
	Branch string // Supports aliases j/fm-3311
}

func Checkout(ctx context.Context, arg CheckoutArgs) error {
	wtCtx, err := GetWorkTreeContext()
	if err != nil {
		return err
	}

	branch := wtCtx.ResolveBranch(arg.Branch)
	slog.Info("checking out", slog.String("branch", branch))

	wt, err := git.FindWorkTree(ctx, branch)
	if err != nil && !errors.Is(err, git.ErrWorkTreeNotFound) {
		return err
	}

	if wt != nil {
		slog.Info("worktree already exists, switching to it")
		return checkoutWorkTree(ctx, wtCtx, wt)
	}

	err = git.Fetch(ctx, "-p")
	if err != nil {
		return err
	}

	// 2. If branch exists on remote, add a new worktree for
	if git.BranchExists(ctx, branch) {
		slog.Info("branch exists on remote, creating new worktree from branch")

		wt, err = git.CreateWorkTreeFromBranch(ctx, wtCtx.Config.WorkTreesDirectory, branch)
		if err != nil {
			return err
		}

		return checkoutWorkTree(ctx, wtCtx, wt)
	}

	mainWt, err := git.FindWorkTree(ctx, "main")
	if err != nil {
		return fmt.Errorf("error finding main worktree: %v", err)
	}

	// Update main worktree
	err = util.InDirectory(mainWt.Path, func() error {
		slog.Info("pulling main")

		return git.Pull(ctx)
	})
	if err != nil {
		return err
	}

	slog.Info("creating new worktree based on main", slog.String("branch", branch))
	wt, err = git.CreateWorkTreeFromNewBranch(ctx, wtCtx.Config.WorkTreesDirectory, branch)
	if err != nil {
		return err
	}

	return checkoutWorkTree(ctx, wtCtx, wt)
}

func checkoutWorkTree(ctx context.Context, wtCtx *WorkTreeContext, wt *git.WorkTree) error {
	slog.Debug("checking out worktree", slog.String("path", wt.Path))

	err := os.Chdir(wt.Path)
	if err != nil {
		return err
	}

	// We don't care if it fails, just want to try and update the branch
	_ = git.Pull(ctx)

	// Copy seed files
	err = wtCtx.SeedWorkTree(wt)
	if err != nil {
		return err
	}

	err = wtCtx.ExecuteAfterCheckoutHooks(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Checked out: %v", wt.Path)

	return nil
}
