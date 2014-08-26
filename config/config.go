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
	File = ".exercism.json"
	// LegacyFile is the name of the original config file.
	// It is a misnomer, since the config was in json, not go.
	LegacyFile = ".exercism.go"
	// Host is the default hostname for fetching problems and submitting exercises.
	// TODO: We need to operate against two hosts (one for problems and one for submissions),
	// or define a proxy that both APIs can go through.
	Host = "http://exercism.io"

	// DirExercises is the default name of the directory for active users.
	DirExercises = "exercism"
)

var (
	errHomeNotFound = errors.New("unable to locate home directory")
)

// Config represents the settings for particular user.
// This defines both the auth for talking to the API, as well as
// where to put problems that get downloaded.
type Config struct {
	APIKey   string `json:"apiKey"`
	Dir      string `json:"exercismDirectory"`
	Hostname string `json:"hostname"`
	home     string // cache user's home directory
	file     string // full path to config file
}

// New returns a new config.
// It will attempt to set defaults where no value is passed in.
func New(key, host, dir string) (*Config, error) {
	c := &Config{
		APIKey:   key,
		Hostname: host,
		Dir:      dir,
	}
	return c.configure()
}

func (c *Config) configure() (*Config, error) {
	c.sanitize()

	if c.Hostname == "" {
		c.Hostname = Host
	}

	dir, err := c.homeDir()
	if err != nil {
		return c, err
	}
	c.file = fmt.Sprintf("%s/%s", dir, File)

	if c.Dir == "" {
		c.Dir = fmt.Sprintf("%s/%s", dir, DirExercises)
	}
	c.Dir = strings.Replace(c.Dir, "~/", fmt.Sprintf("%s/", dir), 1)
	return c, nil
}

// SavePath allows the user to customize the location of the JSON file.
func (c *Config) SavePath(file string) {
	if file != "" {
		c.file = file
	}
}

func (c *Config) File() string {
	return c.file
}

func (c *Config) Write() error {
	renameLegacy()

	// truncates existing file if it exists
	f, err := os.Create(c.file)
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

// ToFile writes a Config to a JSON file.
func (c *Config) ToFile(file string) error {
	file = WithDefaultPath(file)
	f, err := os.Create(file) // truncates existing file if it exists
	if err != nil {
		return err
	}
	defer f.Close()

	c.file = file
	err = c.Encode(f)
	if err != nil {
		return err
	}
	return nil
}

// FromFile loads a Config object from a JSON file.
func FromFile(file string) (*Config, error) {
	file = WithDefaultPath(file)
	f, err := os.Open(file)
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
	c.configure()
	return c, nil
}

// FilePath returns the path to the config file.
func FilePath(file string) (string, error) {
	if file != "" {
		return file, nil
	}

	dir, err := Home()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", dir, File), nil
}

// WithDefaultPath returns the default configuration path if none is provided.
func WithDefaultPath(p string) string {
	if p == "" {
		return fmt.Sprintf("%s/%s", HomeDir(), File)
	}

	return p
}

// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func (c *Config) homeDir() (string, error) {
	if c.home != "" {
		return c.home, nil
	}
	return Home()
}

// Home returns the user's canonical home directory.
// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func Home() (string, error) {
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
		return dir, errHomeNotFound
	}
	return dir, nil
}

// HomeDir returns the user's canonical home directory.
// FIXME: This one will go away. Refactoring in progress.
func HomeDir() string {
	dir, err := Home()
	if err != nil {
		panic("unable to determine the location of your home directory")
	}
	return dir
}

// Demo is a default configuration for unauthenticated users.
func Demo() *Config {
	return &Config{
		Hostname: Host,
		APIKey:   "",
		Dir:      DefaultAssignmentPath(),
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
	c.Dir = strings.TrimSpace(c.Dir)
	c.Hostname = strings.TrimSpace(c.Hostname)
}

// renameLegacy normalizes the default config file name.
// This function will bail silently if any error occurs.
func renameLegacy() {
	dir, err := Home()
	if err != nil {
		return
	}

	legacyPath := filepath.Join(dir, LegacyFile)
	if _, err = os.Stat(legacyPath); err != nil {
		return
	}

	correctPath := filepath.Join(dir, File)
	os.Rename(legacyPath, correctPath)
	return
}
