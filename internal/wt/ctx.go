package wt

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jacobdrury/wt/internal/config"
	"github.com/jacobdrury/wt/internal/git"
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
	// Root is the root path to the `.wt` directory.
	Root string
	// WorkTreesPath is the path to the directory containing the worktrees
	WorkTreesPath string
	// Config is the wt configuration
	Config *config.Config
	// Seed Directory is the directory containing the seed files for new worktrees
	SeedDirectory string
}

func (ctx *WorkTreeContext) initialize() error {
	err := os.Mkdir(ctx.Root, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(ctx.SeedDirectory, 0755)
	if err != nil {
		return err
	}

	return ctx.Config.Save(filepath.Join(ctx.Root, WTConfigFileName))
}

func (ctx *WorkTreeContext) ResolveBranch(val string) string {
	parts := strings.Split(val, ctx.Config.BranchDelimiter)
	for i, part := range parts {
		if alias, ok := ctx.Config.BranchPrefixAliases[config.BranchPrefixAlias(part)]; ok {
			parts[i] = string(alias)
		}
	}

	return strings.Join(parts, ctx.Config.BranchDelimiter)
}

func (ctx *WorkTreeContext) SeedWorkTree(wt *git.WorkTree) error {
	return copy.Copy(ctx.SeedDirectory, wt.Path)
}

func (ctx *WorkTreeContext) ExecuteAfterCheckoutHooks() error {
	for _, hook := range ctx.Config.Hooks.AfterCheckout {
		cmdData := strings.Split(hook, " ")

		cmd := exec.Command(ctx.Config.Hooks.Shell, cmdData...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		print(string(output))
	}

	return nil
}

// CreateContext creates a new default WTContext on the file system in the current
// working directory. The current working directory must be a git repository
// and not have a `.wt` context configured in it.
func CreateContext() error {
	inRepo, err := git.IsGitRepository()
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

	ctx := &WorkTreeContext{
		Root:          wtDir,
		Config:        config.DefaultConfig(),
		SeedDirectory: seedDir,
	}

	return ctx.initialize()
}

// GetContext retrieves the current WTContext from memory.
func GetContext() (*WorkTreeContext, error) {
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

	wtDir := filepath.Join(wd, WTDirectoryName)
	if _, err := os.Stat(wtDir); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotInitialized
		}

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
		Root:          wtDir,
		Config:        cfg,
		SeedDirectory: seedPath,
	}, nil
}
