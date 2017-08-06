package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	cfg.Home = userHome()

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

// Write stores the config to disk.
func (cfg *UserConfig) Write() error {
	cfg.Workspace = cfg.resolve(cfg.Workspace)
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

func (cfg *UserConfig) resolve(path string) string {
	if strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return strings.Replace(path, "~", cfg.Home, 1)
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(cfg.Home, path)
}
