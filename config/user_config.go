package config

import (
	"os"
	"runtime"

	"github.com/spf13/viper"
)

// UserConfig contains user-specific settings.
type UserConfig struct {
	*Config
	Workspace string
	Token     string
	Home      string
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
		Config: New(Dir(), "user"),
	}
}

// SetDefaults ensures that we have proper values where possible.
func (cfg *UserConfig) SetDefaults() {
	if cfg.Home == "" {
		cfg.Home = userHome()
	}
	cfg.Workspace = Resolve(cfg.Workspace, cfg.Home)
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

func userHome() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("USERPROFILE")
		if dir != "" {
			return dir
		}
		dir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if dir != "" {
			return dir
		}
	} else {
		dir = os.Getenv("HOME")
		if dir != "" {
			return dir
		}
	}
	// If all else fails, use the current directory.
	dir, _ = os.Getwd()
	return dir
}
