package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var (
	// DefaultDirName is the default name used for config and workspace directories.
	DefaultDirName string
)

// Configuration lets us inject configuration options into commands.
// Note that we are slowly working towards getting rid of the
// config.Config, config.UserConfig, and config.CLIConfig types.
// Once we do, we can rename this type to Config, and get rid of the
// User and CLI fields.
type Configuration struct {
	OS              string
	Home            string
	Dir             string
	DefaultBaseURL  string
	DefaultDirName  string
	UserViperConfig *viper.Viper
	UserConfig      *UserConfig
	CLIConfig       *CLIConfig
	Persister       Persister
}

// NewConfiguration provides a configuration with default values.
func NewConfiguration() Configuration {
	home := userHome()
	dir := Dir()

	return Configuration{
		OS:             runtime.GOOS,
		Dir:            Dir(),
		Home:           home,
		DefaultBaseURL: defaultBaseURL,
		DefaultDirName: DefaultDirName,
		Persister:      FilePersister{Dir: dir},
	}
}

// SetDefaultDirName configures the default directory name based on the name of the binary.
func SetDefaultDirName(binaryName string) {
	DefaultDirName = strings.Replace(filepath.Base(binaryName), ".exe", "", 1)
}

// Dir is the configured config home directory.
// All the cli-related config files live in this directory.
func Dir() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir != "" {
			return filepath.Join(dir, DefaultDirName)
		}
	} else {
		dir := os.Getenv("EXERCISM_CONFIG_HOME")
		if dir != "" {
			return dir
		}
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = filepath.Join(os.Getenv("HOME"), ".config")
		}
		if dir != "" {
			return filepath.Join(dir, DefaultDirName)
		}
	}
	// If all else fails, use the current directory.
	dir, _ = os.Getwd()
	return dir
}

func userHome() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("USERPROFILE")
		if dir != "" {
			return dir
		}
		dir = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
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

// DefaultWorkspaceDir provides a sensible default for the Exercism workspace.
// The default is different depending on the platform, in order to best match
// the conventions for that platform.
// It places the directory in the user's home path.
func DefaultWorkspaceDir(cfg Configuration) string {
	dir := cfg.DefaultDirName
	if cfg.OS != "linux" {
		dir = strings.Title(dir)
	}
	return filepath.Join(cfg.Home, dir)
}

// Save persists a viper config of the base name.
func (c Configuration) Save(basename string) error {
	return c.Persister.Save(c.UserViperConfig, basename)
}
