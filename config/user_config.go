package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	defaultBaseURL = "https://v2.exercism.io/api/v1"
)

// UserConfig contains user-specific settings.
type UserConfig struct {
	*Config
	Workspace  string
	Token      string
	Home       string
	APIBaseURL string
	settings   Configuration
}

// NewUserConfig loads a user configuration if it exists.
func NewUserConfig() (*UserConfig, error) {
	cfg := NewEmptyUserConfig()

	if err := cfg.Load(viper.New()); err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewEmptyUserConfig creates a user configuration without loading it.
func NewEmptyUserConfig() *UserConfig {
	return &UserConfig{
		Config:   New(Dir(), "user"),
		settings: NewConfiguration(),
	}
}

// SetDefaults ensures that we have proper values where possible.
func (cfg *UserConfig) SetDefaults() {
	if cfg.Home == "" {
		cfg.Home = userHome()
	}
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = defaultBaseURL
	}
	if cfg.Workspace == "" {
		dir := DefaultWorkspaceDir(cfg.settings)

		_, err := os.Stat(dir)
		// Sorry about the double negative.
		if !os.IsNotExist(err) {
			dir = fmt.Sprintf("%s-1", dir)
		}
		cfg.Workspace = dir
	}
}

// Write stores the config to disk.
func (cfg *UserConfig) Write() error {
	cfg.SetDefaults()
	return Write(cfg)
}

// Load reads a viper configuration into the config.
func (cfg *UserConfig) Load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}
