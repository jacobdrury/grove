package grove

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jacobdrury/grove/internal/config"
	"github.com/jacobdrury/grove/internal/git"
)

var (
	instance *Grove

	ErrAlreadyInitialized    = errors.New("already initialized")
	ErrNotAGitRepository     = errors.New("not a git repository")
	ErrNotInitialized        = errors.New("not initialized")
	ErrNotLoaded             = errors.New("not loaded, call wtcontext.Load() first")
	ErrConfigNotFound        = errors.New("config not found")
	ErrSeedDirectoryNotFound = errors.New("seed directory not found")
)

const (
	GroveDirectoryName string = ".grove"
	SeedDirectoryName  string = "seed"
	ConfigFileName     string = "config.yaml"
)

type Grove struct {
	// Config is the grove configuration
	Config *config.Config

	// RepositoryPath is the root directory in which the `.grove` directory is stored
	RepositoryPath string
	// GrovePath is the root path to the `.grove` directory.
	GrovePath string
	// WorkTreesPath is the path to the directory containing the worktrees
	WorkTreesPath string
	// Seed Directory is the directory containing the seed files for new worktrees
	SeedPath string
}

func (grove *Grove) persist() error {
	err := os.Mkdir(grove.GrovePath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(grove.SeedPath, 0755)
	if err != nil {
		return err
	}

	return grove.Config.Save(filepath.Join(grove.GrovePath, ConfigFileName))
}

// New creates a new default Grove on the file system in the current
// working directory. The current working directory must be a git repository
// and not have a `.grove` directory configured in it.
func New(ctx context.Context) (*Grove, error) {
	inRepo, err := git.IsGitRepository(ctx)
	if err != nil {
		return nil, err
	}

	if !inRepo {
		return nil, ErrNotAGitRepository
	}

	err = Load()
	if err != nil {
		if !errors.Is(err, ErrNotInitialized) {
			return nil, err
		}
	}

	if err == nil {
		return nil, ErrAlreadyInitialized
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wtDir := filepath.Join(wd, GroveDirectoryName)
	seedDir := filepath.Join(wtDir, SeedDirectoryName)

	grove := &Grove{
		GrovePath: wtDir,
		Config:    config.DefaultConfig(),
		SeedPath:  seedDir,
	}

	err = grove.persist()
	if err != nil {
		return nil, err
	}

	return grove, nil
}

// GetInstance retrieves the current Grove instance from memory.
func GetInstance() (*Grove, error) {
	if instance == nil {
		return nil, ErrNotLoaded
	}

	return instance, nil
}

// Load loads the Grove instance from the file-system into memory.
func Load() error {
	ctx, err := load()
	if err != nil {
		return err
	}

	instance = ctx
	return nil
}

func load() (*Grove, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	groveDir, err := locateGroveDir(wd)
	if err != nil {
		return nil, err
	}

	seedPath := filepath.Join(groveDir, SeedDirectoryName)
	if _, err := os.Stat(seedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSeedDirectoryNotFound
		}

		return nil, err
	}

	cfgPath := filepath.Join(groveDir, ConfigFileName)
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

	return &Grove{
		RepositoryPath: filepath.Dir(groveDir),
		GrovePath:      groveDir,
		Config:         cfg,
		SeedPath:       seedPath,
	}, nil
}

func locateGroveDir(startPath string) (string, error) {
	dir := startPath
	for {
		wtPath := filepath.Join(dir, GroveDirectoryName)
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
