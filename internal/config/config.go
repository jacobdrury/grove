package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type (
	BranchPrefixAlias string
	BranchPrefix      string
)

type Hooks struct {
	Shell         string   `yaml:"shell"`
	AfterCheckout []string `yaml:"after-checkout"`
}

type BranchResolver struct {
	BranchDelimiter     string                             `yaml:"branch-delimiter"`
	BranchPrefixAliases map[BranchPrefixAlias]BranchPrefix `yaml:"prefix-aliases"`
}

type Config struct {
	WorkTreesDirectory string         `yaml:"worktrees-directory"`
	BranchResolver     BranchResolver `yaml:"branch-resolver"`
	Hooks              Hooks          `yaml:"hooks"`
}

func DefaultConfig() *Config {
	var defaultShell string
	switch runtime.GOOS {
	case "windows":
		defaultShell = os.Getenv("ComSpec")
		if defaultShell == "" {
			defaultShell = "C:\\Windows\\system32\\cmd.exe"
		}
	default:
		defaultShell = os.Getenv("SHELL")
		if defaultShell == "" {
			defaultShell = "/bin/sh"
		}
	}

	return &Config{
		WorkTreesDirectory: "./worktrees",
		BranchResolver: BranchResolver{
			BranchPrefixAliases: map[BranchPrefixAlias]BranchPrefix{},
			BranchDelimiter:     "/",
		},
		Hooks: Hooks{
			Shell:         defaultShell,
			AfterCheckout: []string{},
		},
	}
}

// Save saves the configuration to the specified path.
func (c *Config) Save(path string) error {
	marshaled, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, marshaled, 0644)
}

// Load loads the config at the specified path into memory.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
