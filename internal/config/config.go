package config

import (
	"os"
	"runtime"
	"strings"

	"github.com/samber/lo"
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
	BranchPrefixAliases map[BranchPrefixAlias]BranchPrefix `yaml:"prefix-aliases"`
	BranchDelimiter     string                             `yaml:"branch-delimiter"`
}

func (b *BranchResolver) Resolve(val string, branches []string) string {
	originalParts := strings.Split(val, b.BranchDelimiter)
	parts := originalParts

	for i, part := range parts {
		// Expand aliases
		if i != len(parts)-1 {
			if alias, ok := b.BranchPrefixAliases[BranchPrefixAlias(part)]; ok {
				parts[i] = string(alias)
			}

			continue
		}

		// Resolve by exact match
		if branch, ok := lo.Find(branches, func(branch string) bool {
			return branch == strings.Join(parts, b.BranchDelimiter)
		}); ok {
			return branch
		}

		// Resolve slug by prefix match
		for _, branch := range branches {
			branchParts := strings.Split(branch, b.BranchDelimiter)
			slug := branchParts[len(branchParts)-1]
			prefix := strings.Join(branchParts[:len(branchParts)-1], b.BranchDelimiter)
			resolvedPrefix := strings.Join(parts[:len(parts)-1], b.BranchDelimiter)

			if strings.HasPrefix(slug, part) && prefix == resolvedPrefix {
				parts[i] = slug
				break
			}
		}
	}

	return strings.Join(parts, b.BranchDelimiter)
}

type Config struct {
	WorkTreesDirectory string         `yaml:"worktrees-directory"`
	BranchResolver     BranchResolver `yaml:"branch-resolver"`
	Hooks              Hooks          `yaml:"hooks"`
}

// Save saves the configuration to the specified path.
func (c *Config) Save(path string) error {
	marshaled, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, marshaled, 0644)
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
