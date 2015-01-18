package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	fileEnvKey = "EXERCISM_CONFIG_FILE"
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
	// Make this non-exported when handlers.Login is deleted.
	DirExercises = "exercism"
)

var (
	errHomeNotFound = errors.New("unable to locate home directory")
)

// Config represents the settings for particular user.
// This defines both the auth for talking to the API, as well as
// where to put problems that get downloaded.
type Config struct {
	APIKey string `json:"apiKey"`
	Dir    string `json:"dir"`
	API    string `json:"api"`
	XAPI   string `json:"xapi"`
	File   string `json:"-"` // full path to config file
	home   string // cache user's home directory
}

// Home returns the user's canonical home directory.
// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func Home() (string, error) {
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
		return dir, errHomeNotFound
	}
	return dir, nil
}

func New(path string) (*Config, error) {
	c := &Config{}
	err := c.load(path, os.Getenv(fileEnvKey))
	return c, err
}

// Update sets new values where given.
func (c *Config) Update(key, host, dir, xapi string) {
	key = strings.TrimSpace(key)
	if key != "" {
		c.APIKey = key
	}

	host = strings.TrimSpace(host)
	if host != "" {
		c.API = host
	}

	dir = strings.TrimSpace(dir)
	if dir != "" {
		c.Dir = dir
	}

	xapi = strings.TrimSpace(xapi)
	if xapi != "" {
		c.XAPI = xapi
	}
}

// Write saves the config as JSON.
func (c *Config) Write() error {
	// truncates existing file if it exists
	f, err := os.Create(c.File)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(c)
}

func (c *Config) load(argPath, envPath string) error {
	path, err := c.resolvePath(argPath, envPath)
	if err != nil {
		return err
	}
	c.File = path

	if err := c.read(); err != nil {
		return err
	}

	// in case people manually update the config file
	// with weird formatting
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.Dir = strings.TrimSpace(c.Dir)
	c.API = strings.TrimSpace(c.API)
	c.XAPI = strings.TrimSpace(c.XAPI)

	return c.setDefaults()
}

func (c *Config) read() error {
	if _, err := os.Stat(c.File); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	f, err := os.Open(c.File)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	return d.Decode(&c)
}

// IsAuthenticated returns true if the config contains an API key.
// This does not check whether or not that key is valid.
func (c *Config) IsAuthenticated() bool {
	return c.APIKey != ""
}

// homeDir caches the lookup of the user's home directory.
func (c *Config) homeDir() (string, error) {
	if c.home != "" {
		return c.home, nil
	}
	return Home()
}

func (c *Config) resolvePath(argPath, envPath string) (string, error) {
	path := argPath
	if path == "" {
		path = envPath
	}
	if path == "" {
		path = filepath.Join("~", File)
	}
	h, err := c.homeDir()
	if err != nil {
		return "", err
	}
	return strings.Replace(path, "~", h, 1), nil
}

func (c *Config) setDefaults() error {
	if c.API == "" {
		c.API = hostAPI
	}

	if c.XAPI == "" {
		c.XAPI = hostXAPI
	}

	h, err := c.homeDir()
	if err != nil {
		return err
	}

	if c.Dir == "" {
		c.Dir = filepath.Join(h, DirExercises)
	}
	c.Dir = strings.Replace(c.Dir, "~", h, 1)

  return nil
}
