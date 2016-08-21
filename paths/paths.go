package paths

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// File is the default name of the JSON file where the config written.
	// The user can pass an alternate filename when using the CLI.
	File = "exercism.json"
	// DirExercises is the default name of the directory for active users.
	// Make this non-exported when handlers.Login is deleted.
	DirExercises = "exercism"
)

var (
	// Home by default will contact the location of your home directory.
	Home string
	// ConfigHome will contain $XDG_CONFIG_HOME if it is set or default config home directory.
	ConfigHome string
	// DefaultConfig will contain default path to config, according to Home
	DefaultConfig string

	errHomeNotFound = errors.New("unable to locate home directory")
)

func init() {
	var err error
	Home, err = findHome()
	if err != nil {
		panic(err)
	}
	ConfigHome = os.Getenv("XDG_CONFIG_HOME")
	if ConfigHome == "" {
		ConfigHome = filepath.Join(Home, ".config")
	}
	DefaultConfig = filepath.Join(Home, "." + File)
}

// Config will return the correct input path given any input.
// Blank input will return the default configuration location based
// on ConfigHome.
// Non-blank input will expand home to be an absolute path.
// If the target is known to be a directory, the config filename
// will be appended.
func Config(path string) string {
	if path == "" {
		return filepath.Join(ConfigHome, File)
	}
	expandedPath := expandPath(path)
	if IsDir(path) {
		expandedPath = filepath.Join(expandedPath, File)
	}
	return expandedPath
}

// Exercises will return the correct exercises path given any input.
// Blank input will return the default location for exercises.
// Non-blank input will expand home to be an absolute path.
func Exercises(path string) string {
	if path == "" {
		return filepath.Join(Home, DirExercises)
	}
	return expandPath(path)
}

// IsDir determines whether the given path is a valid directory path.
func IsDir(path string) bool {
	fi, _ := os.Stat(path)
	return fi != nil && fi.IsDir()
}

func expandPath(path string) string {
	return makeAbsolute(expandHome(strings.TrimSpace(path)))
}

func findHome() (string, error) {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("USERPROFILE")
		if dir == "" {
			dir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	} else {
		dir = os.Getenv("HOME")
	}

	if dir == "" {
		return "", errHomeNotFound
	}

	return dir, nil
}

func makeAbsolute(path string) string {
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return filepath.Join(wd, path)
	}
	return path
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return strings.Replace(path, "~", Home, 1)
	}
	return path
}
