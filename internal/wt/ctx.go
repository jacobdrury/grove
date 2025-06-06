package wt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobdrury/wt/internal/config"
	"github.com/jacobdrury/wt/internal/git"
	"github.com/jacobdrury/wt/internal/util"
	"github.com/otiai10/copy"
)

var (
	wtContext *WorkTreeContext

	ErrAlreadyInitialized    = errors.New("already initialized")
	ErrNotAGitRepository     = errors.New("not a git repository")
	ErrNotInitialized        = errors.New("not initialized")
	ErrNotLoaded             = errors.New("not loaded, call wtcontext.Load() first")
	ErrConfigNotFound        = errors.New("config not found")
	ErrSeedDirectoryNotFound = errors.New("seed directory not found")
)

const (
	WTDirectoryName     string = ".wt"
	WTSeedDirectoryName string = "seed"
	WTConfigFileName    string = "config.yaml"
)

type WorkTreeContext struct {
	// Root is the root directory in which the `.wt` directory is stored
	Root string
	// WTDirectory is the root path to the `.wt` directory.
	WTDirectory string
	// WorkTreesPath is the path to the directory containing the worktrees
	WorkTreesPath string
	// Config is the wt configuration
	Config *config.Config
	// Seed Directory is the directory containing the seed files for new worktrees
	SeedDirectory string
}

func (wtCtx *WorkTreeContext) initialize() error {
	err := os.Mkdir(wtCtx.WTDirectory, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(wtCtx.SeedDirectory, 0755)
	if err != nil {
		return err
	}

	return wtCtx.Config.Save(filepath.Join(wtCtx.WTDirectory, WTConfigFileName))
}

func (wtCtx *WorkTreeContext) ResolveBranch(val string) string {
	parts := strings.Split(val, wtCtx.Config.BranchDelimiter)
	for i, part := range parts {
		if alias, ok := wtCtx.Config.BranchPrefixAliases[config.BranchPrefixAlias(part)]; ok {
			parts[i] = string(alias)
		}
	}

	return strings.Join(parts, wtCtx.Config.BranchDelimiter)
}

func (wtCtx *WorkTreeContext) SeedWorkTree(wt *git.WorkTree) error {
	slog.Debug("seeding worktree", slog.String("workTreePath", wt.Path), slog.String("seedDirectory", wtCtx.SeedDirectory))
	return copy.Copy(wtCtx.SeedDirectory, wt.Path)
}

func (wtCtx *WorkTreeContext) ExecuteAfterCheckoutHooks(ctx context.Context) error {
	slog.Debug("executing after checkout hooks", slog.Int("numberOfHooks", len(wtCtx.Config.Hooks.AfterCheckout)))
	for _, hook := range wtCtx.Config.Hooks.AfterCheckout {
		slog.Info("executing hook", slog.String("hook", hook))

		err := util.ExecShellCmd(ctx, wtCtx.Config.Hooks.Shell, hook)
		if err != nil {
			return fmt.Errorf("error executing hook %s: %v", hook, err)
		}
	}

	slog.Debug("after checkout hooks executed")

	return nil
}

// CreateContext creates a new default WTContext on the file system in the current
// working directory. The current working directory must be a git repository
// and not have a `.wt` context configured in it.
func CreateContext(ctx context.Context) error {
	inRepo, err := git.IsGitRepository(ctx)
	if err != nil {
		return err
	}

	if !inRepo {
		return ErrNotAGitRepository
	}

	err = LoadContext()
	if err != nil {
		if !errors.Is(err, ErrNotInitialized) {
			return err
		}
	}

	if err == nil {
		return ErrAlreadyInitialized
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	wtDir := filepath.Join(wd, WTDirectoryName)
	seedDir := filepath.Join(wtDir, WTSeedDirectoryName)

	wtCtx := &WorkTreeContext{
		WTDirectory:   wtDir,
		Config:        config.DefaultConfig(),
		SeedDirectory: seedDir,
	}

	return wtCtx.initialize()
}

// GetWorkTreeContext retrieves the current WTContext from memory.
func GetWorkTreeContext() (*WorkTreeContext, error) {
	if wtContext == nil {
		return nil, ErrNotLoaded
	}

	return wtContext, nil
}

// LoadContext loads the WTContext from the file-system into memory.
func LoadContext() error {
	ctx, err := loadContext()
	if err != nil {
		return err
	}

	wtContext = ctx
	return nil
}

func loadContext() (*WorkTreeContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wtDir, err := locateWtDir(wd)
	if err != nil {
		return nil, err
	}

	seedPath := filepath.Join(wtDir, WTSeedDirectoryName)
	if _, err := os.Stat(seedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSeedDirectoryNotFound
		}

		return nil, err
	}

	cfgPath := filepath.Join(wtDir, WTConfigFileName)
	if _, err := os.Stat(cfgPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}

		return nil, err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return &WorkTreeContext{
		Root:          filepath.Dir(wtDir),
		WTDirectory:   wtDir,
		Config:        cfg,
		SeedDirectory: seedPath,
	}, nil
}

func locateWtDir(startPath string) (string, error) {
	dir := startPath
	for {
		wtPath := filepath.Join(dir, WTDirectoryName)
		info, err := os.Stat(wtPath)
		if err == nil && info.IsDir() {
			return wtPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", ErrNotInitialized
}
