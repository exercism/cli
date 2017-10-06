package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Dir is the configured config home directory.
// All the cli-related config files live in this directory.
func Dir() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir != "" {
			return filepath.Join(dir, "exercism")
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
			return filepath.Join(dir, "exercism")
		}
	}
	// If all else fails, use the current directory.
	dir, _ = os.Getwd()
	return dir
}
