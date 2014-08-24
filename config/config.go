package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// File is the default name of the JSON file where the config written.
	// The user can pass an alternate filename when using the CLI.
	File       = ".exercism.json"
	LegacyFile = ".exercism.go"
	// Host is the default hostname for fetching problems and submitting exercises.
	// TODO: We need to operate against two hosts (one for problems and one for submissions),
	// or define a proxy that both APIs can go through.
	Host = "http://exercism.io"
	// DemoDirname is the default directory to download problems to.
	DemoDirname = "exercism-demo"

	// AssignmentDirname is the default name of the directory for active users.
	AssignmentDirname = "exercism"
)

// Config represents the settings for particular user.
// This defines both the auth for talking to the API, as well as
// where to put problems that get downloaded.
type Config struct {
	APIKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
	Hostname          string `json:"hostname"`
}

// ToFile writes a Config to a JSON file.
func (c Config) ToFile(path string) error {
	if path == "" {
		path = WithDefaultPath(path)
		err := normalizeFilename(HomeDir())
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path) // truncates existing file if it exists
	if err != nil {
		return err
	}
	defer f.Close()

	err = c.Encode(f)
	if err != nil {
		return err
	}
	return nil
}

// FromFile loads a Config object from a JSON file.
func FromFile(path string) (*Config, error) {
	if path == "" {
		path = WithDefaultPath(path)
		err := normalizeFilename(HomeDir())
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

// Encode writes a Config into JSON format.
func (c *Config) Encode(w io.Writer) error {
	c.sanitize()
	e := json.NewEncoder(w)
	return e.Encode(c)
}

// Decode loads a Config from JSON format.
func Decode(r io.Reader) (*Config, error) {
	d := json.NewDecoder(r)
	var c *Config
	err := d.Decode(&c)
	if err != nil {
		return c, err
	}
	c.sanitize()

	return c, err
}

// WithDefaultPath returns the default configuration path if none is provided.
func WithDefaultPath(p string) string {
	if p == "" {
		return Filename(HomeDir())
	}

	return p
}

var homeDir string

// HomeDir returns the user's canonical home directory.
// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func HomeDir() string {
	if homeDir != "" {
		return homeDir
	}

	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE")
		}
	} else {
		homeDir = os.Getenv("HOME")
	}

	// TODO should we fall back to the CWD instead ?
	if homeDir == "" {
		panic("unable to determine the location of your home directory")
	}

	return homeDir

}

// Filename is the name of the JSON file containing the user's config.
func Filename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, File)
}

// Demo is a default configuration for unauthenticated users.
func Demo() *Config {
	return &Config{
		Hostname:          Host,
		APIKey:            "",
		ExercismDirectory: demoDirectory(),
	}
}

// DefaultAssignmentPath returns the absolute path of the default exercism directory
func DefaultAssignmentPath() string {
	return filepath.Join(HomeDir(), AssignmentDirname)
}

// ReplaceTilde replaces the short-hand home path with the absolute path.
func ReplaceTilde(oldPath string) string {
	return strings.Replace(oldPath, "~/", HomeDir()+"/", 1)
}

func normalizeFilename(path string) error {
	var err error

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("expected path to be a directory")
	}

	currentPath := filepath.Join(path, File)
	oldPath := filepath.Join(path, LegacyFile)

	_, err = os.Stat(currentPath)
	// Do nothing nil means we already have a current config file
	if err == nil {
		return nil
	}
	// return any error unless the error is because the file is missing
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Do nothing if we have no old file to rename
	_, err = os.Stat(oldPath)
	if os.IsNotExist(err) {
		return nil
	}

	err = os.Rename(oldPath, currentPath)
	if err != nil {
		return err
	}
	fmt.Printf("renamed %s to %s\n", oldPath, currentPath)

	return nil
}

func demoDirectory() string {
	return filepath.Join(HomeDir(), DemoDirname)
}

func (c *Config) sanitize() {
	c.APIKey = sanitizeField(c.APIKey)
	c.ExercismDirectory = sanitizeField(c.ExercismDirectory)
	c.Hostname = sanitizeField(c.Hostname)
}

func sanitizeField(v string) string {
	return strings.TrimSpace(v)
}
