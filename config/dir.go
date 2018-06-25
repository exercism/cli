package config

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// DefaultDirName is the default name used for config and workspace directories.
	DefaultDirName string
)

// SetDefaultDirName configures the default directory name based on the name of the binary.
func SetDefaultDirName(binaryName string) {
	DefaultDirName = strings.Replace(path.Base(binaryName), ".exe", "", 1)
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
