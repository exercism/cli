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
	File = ".exercism.json"
	// DirExercises is the default name of the directory for active users.
	// Make this non-exported when handlers.Login is deleted.
	DirExercises = "exercism"
)

var (
	// Home by default will contact the location of your home directory
	Home            string
	errHomeNotFound = errors.New("unable to locate home directory")
)

func init() {
	// on startup set default values for Home and Config
	Recalculate()
}

// Config will return the correct input path given any input.
//  blank input will return the default location for configuration
//  non-blank input will have:
//  * its home expanded and will be expanded to an absolute path
//  * the config file appended if we know the target is a directory
func Config(path string) string {
	if path == "" {
		return filepath.Join(Home, File)
	}

	expandedPath := expandPath(path)
	if IsDir(path) {
		expandedPath = filepath.Join(expandedPath, File)
	}
	return expandedPath
}

// Exercises will return the correct exercises path given any input.
//  blank input will return the default location for exercises
//  non-blank input will have:
//  * its home expanded and will be expanded to an absolute path
//  * the config file appended if we know the target is a directory
func Exercises(path string) string {
	if path == "" {
		return filepath.Join(Home, DirExercises)
	}
	return expandPath(path)
}

// Recalculate sets exercism paths based on Home
func Recalculate() {
	if Home == "" {
		home, err := findHome()
		if err != nil {
			panic(err)
		}
		Home = home
	}
}

// IsDir tells us if the path is valid and is a directory
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
	if path[:2] == "~"+string(os.PathSeparator) {
		return strings.Replace(path, "~", Home, 1)
	}
	return path
}
