package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var (
	defaultBaseURL = "https://api.exercism.org/v1"

	// DefaultDirName is the default name used for config and workspace directories.
	DefaultDirName string
)

// Config lets us inject configuration options into commands.
type Config struct {
	OS              string
	Home            string
	Dir             string
	DefaultBaseURL  string
	DefaultDirName  string
	UserViperConfig *viper.Viper
	Persister       Persister
}

// NewConfig provides a configuration with default values.
func NewConfig() Config {
	home := userHome()
	dir := Dir()

	return Config{
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
func DefaultWorkspaceDir(cfg Config) string {
	dir := cfg.DefaultDirName
	if cfg.OS != "linux" {
		dir = strings.Title(dir)
	}
	return filepath.Join(cfg.Home, dir)
}

// Save persists a viper config of the base name.
func (c Config) Save(basename string) error {
	return c.Persister.Save(c.UserViperConfig, basename)
}

// InferSiteURL guesses what the website URL is.
// The basis for the guess is which API we're submitting to.
func InferSiteURL(apiURL string) string {
	if apiURL == "" {
		apiURL = defaultBaseURL
	}
	if apiURL == "https://api.exercism.org/v1" {
		return "https://exercism.org"
	}
	re := regexp.MustCompile("^(https?://[^/]*).*")
	return re.ReplaceAllString(apiURL, "$1")
}

// SettingsURL provides a link to where the user can find their API token.
func SettingsURL(apiURL string) string {
	return fmt.Sprintf("%s%s", InferSiteURL(apiURL), "/my/settings")
}
