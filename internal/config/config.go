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

type Config struct {
	WorkTreesDirectory  string                             `yaml:"worktrees-directory"`
	BranchPrefixAliases map[BranchPrefixAlias]BranchPrefix `yaml:"prefix-aliases"`
	BranchDelimiter     string                             `yaml:"branch-delimiter"`
	Hooks               Hooks                              `yaml:"hooks"`

	// BranchSlugFormatPattern *regexp.Regexp                     `yaml:"slug-format"` // FUTURE: :(
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
		defaultShell = "powershell"
	default:
		defaultShell = "sh"
	}

	return &Config{
		WorkTreesDirectory:  "./worktrees",
		BranchPrefixAliases: map[BranchPrefixAlias]BranchPrefix{},
		BranchDelimiter:     "/",
		Hooks: Hooks{
			Shell:         defaultShell,
			AfterCheckout: []string{},
		},

		// BranchSlugFormatPattern: regexp.MustCompile(`^.*$`), // FUTURE: :(
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
