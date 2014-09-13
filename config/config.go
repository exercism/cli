package config

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// hostAPI is the endpoint to submit solutions to, and to get personalized data
	hostAPI = "http://exercism.io"
	// hostXAPI is the endpoint to fetch problems from
	hostXAPI = "http://x.exercism.io"

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
	APIKey       string `json:"apiKey"`
	Dir          string `json:"exercismDirectory"`
	Hostname     string `json:"hostname"`
	ProblemsHost string `json:"problemsHost"`
	home         string // cache user's home directory
	file         string // full path to config file
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
		c.Hostname = hostAPI
	}

	if c.ProblemsHost == "" {
		c.ProblemsHost = hostXAPI
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

// File represents the path to the config file.
func (c *Config) File() string {
	return c.file
}

// Write() saves the config as JSON.
func (c *Config) Write() error {
	renameLegacy()

	// truncates existing file if it exists
	f, err := os.Create(c.file)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(c)
}

// Read loads the config from the stored JSON file.
func Read(file string) (*Config, error) {
	renameLegacy()

	file, err := FilePath(file)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return New("", "", "")
		}
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c *Config
	d := json.NewDecoder(f)
	err = d.Decode(&c)
	if err != nil {
		return c, err
	}
	c.SavePath(file)
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

// IsAuthenticated returns true if the config contains an API key.
// This does not check whether or not that key is valid.
func (c *Config) IsAuthenticated() bool {
	return c.APIKey != ""
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

func (c *Config) sanitize() {
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.Dir = strings.TrimSpace(c.Dir)
	c.Hostname = strings.TrimSpace(c.Hostname)
	c.ProblemsHost = strings.TrimSpace(c.ProblemsHost)
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
