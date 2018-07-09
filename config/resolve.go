package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Resolve cleans up filesystem paths.
func Resolve(path, home string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~/", "", 1)
		return filepath.Join(home, path)
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	// if using "/dir" on Windows
	if strings.HasPrefix(path, "/") {
		return filepath.Join(home, filepath.Clean(path))
	}
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}
	return filepath.Join(cwd, path)
}
