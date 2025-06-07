package grove

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jacobdrury/grove/internal/config"
	"github.com/jacobdrury/grove/internal/git"
	"github.com/jacobdrury/grove/internal/util"
	"github.com/otiai10/copy"
	"github.com/samber/lo"
)

type CheckoutArgs struct {
	Branch string // Supports aliases j/fm-3311
}

func (grove *Grove) Checkout(ctx context.Context, arg CheckoutArgs) error {
	return util.InDirectory(grove.RepositoryPath, func() error {
		branches, err := git.ListBranches(ctx)
		if err != nil {
			return err
		}

		name := arg.Branch
		branch := grove.resolveBranch(name, branches)
		slog.Info("checking out", slog.String("branch", branch))

		wt, err := git.FindWorkTree(ctx, branch)
		if err != nil && !errors.Is(err, git.ErrWorkTreeNotFound) {
			return err
		}

		if wt != nil {
			slog.Info("worktree already exists, switching to it")
			return checkoutWorkTree(ctx, grove, wt)
		}

		err = git.Fetch(ctx, "-p")
		if err != nil {
			return err
		}

		// 2. If branch exists on remote, add a new worktree for
		if git.BranchExists(ctx, branch) {
			slog.Info("branch exists on remote, creating new worktree from branch")

			wt, err = git.CreateWorkTreeFromBranch(ctx, grove.Config.WorkTreesDirectory, name, branch)
			if err != nil {
				return err
			}

			return checkoutWorkTree(ctx, grove, wt)
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
		wt, err = git.CreateWorkTreeFromNewBranch(ctx, grove.Config.WorkTreesDirectory, name, branch)
		if err != nil {
			return err
		}

		return checkoutWorkTree(ctx, grove, wt)
	})
}

func (grove *Grove) seedWorkTree(wt *git.WorkTree) error {
	slog.Debug("seeding worktree", slog.String("workTreePath", wt.Path), slog.String("seedDirectory", grove.SeedPath))

	return copy.Copy(grove.SeedPath, wt.Path)
}

func (grove *Grove) resolveBranch(val string, branches []string) string {
	br := grove.Config.BranchResolver

	originalParts := strings.Split(val, br.BranchDelimiter)
	parts := originalParts

	for i, part := range parts {
		// Expand aliases
		if i != len(parts)-1 {
			if alias, ok := br.BranchPrefixAliases[config.BranchPrefixAlias(part)]; ok {
				parts[i] = string(alias)
			}

			continue
		}

		// Resolve by exact match
		if branch, ok := lo.Find(branches, func(branch string) bool {
			return branch == strings.Join(parts, br.BranchDelimiter)
		}); ok {
			return branch
		}

		// Resolve slug by prefix match
		for _, branch := range branches {
			branchParts := strings.Split(branch, br.BranchDelimiter)
			slug := branchParts[len(branchParts)-1]
			prefix := strings.Join(branchParts[:len(branchParts)-1], br.BranchDelimiter)
			resolvedPrefix := strings.Join(parts[:len(parts)-1], br.BranchDelimiter)

			if strings.HasPrefix(slug, part) && prefix == resolvedPrefix {
				parts[i] = slug
				break
			}
		}
	}

	return strings.Join(parts, br.BranchDelimiter)
}

func checkoutWorkTree(ctx context.Context, grove *Grove, wt *git.WorkTree) error {
	slog.Debug("checking out worktree", slog.String("path", wt.Path))

	err := os.Chdir(wt.Path)
	if err != nil {
		return err
	}

	// We don't care if it fails, just want to try and update the branch
	_ = git.Pull(ctx)

	// Copy seed files
	err = grove.seedWorkTree(wt)
	if err != nil {
		return err
	}

	err = grove.executeAfterCheckoutHooks(ctx)
	if err != nil {
		return err
	}

	slog.Info("checked out worktree", slog.String("path", wt.Path))

	return nil
}
