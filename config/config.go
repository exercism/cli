package config

import (
	"encoding/json"
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

	// DirExercises is the default name of the directory for active users.
	DirExercises = "exercism"
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
	path = WithDefaultPath(path)
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
	path = WithDefaultPath(path)
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
		return fmt.Sprintf("%s/%s", HomeDir(), File)
	}

	return p
}

// HomeDir returns the user's canonical home directory.
// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func HomeDir() string {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if dir == "" {
			dir = os.Getenv("USERPROFILE")
		}
	} else {
		dir = os.Getenv("HOME")
	}

	if dir == "" {
		panic("unable to determine the location of your home directory")
	}

	return dir
}

// Demo is a default configuration for unauthenticated users.
func Demo() *Config {
	return &Config{
		Hostname:          Host,
		APIKey:            "",
		ExercismDirectory: DefaultAssignmentPath(),
	}
}

// DefaultAssignmentPath returns the absolute path of the default exercism directory
func DefaultAssignmentPath() string {
	return filepath.Join(HomeDir(), DirExercises)
}

// ReplaceTilde replaces the short-hand home path with the absolute path.
func ReplaceTilde(path string) string {
	return strings.Replace(path, "~/", HomeDir()+"/", 1)
}

func (c *Config) sanitize() {
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.ExercismDirectory = strings.TrimSpace(c.ExercismDirectory)
	c.Hostname = strings.TrimSpace(c.Hostname)
}
