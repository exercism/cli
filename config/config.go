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

	// deprecated, get rid of them when nobody uses 1.7.0 anymore
	ExercismDirectory string `json:"exercismDirectory,omitempty"`
	Hostname          string `json:"hostname,omitempty"`
	ProblemsHost      string `json:"problemsHost,omitempty"`
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

// Read loads the config from the stored JSON file.
func Read(file string) (*Config, error) {
	c := &Config{}
	err := c.Read(file)
	return c, err
}

// Update sets new values where given.
func (c *Config) Update(key, host, dir, xapi string) {
	if key != "" {
		c.APIKey = key
	}

	if host != "" {
		c.API = host
	}

	if dir != "" {
		c.Dir = dir
	}

	if xapi != "" {
		c.XAPI = xapi
	}

	c.configure()
}

// Expand takes inputs for a config file location and builds an absolute path.
func Expand(path, env, home string) string {
	if path == "" {
		path = env
	}

	if path != "" && path[0] == '~' {
		path = strings.Replace(path, "~/", fmt.Sprintf("%s/", home), 1)
	}

	if path == "" {
		path = filepath.Join(home, File)
	}

	return path
}

// Read loads the config from the stored JSON file.
func (c *Config) Read(file string) error {
	renameLegacy()

	home, err := c.homeDir()
	if err != nil {
		return err
	}

	c.File = Expand(file, os.Getenv(fileEnvKey), home)

	if _, err := os.Stat(c.File); err != nil {
		if os.IsNotExist(err) {
			c.configure()
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
	err = d.Decode(&c)
	if err != nil {
		return err
	}
	c.configure()
	return nil
}

// SavePath allows the user to customize the location of the JSON file.
func (c *Config) SavePath(file string) {
	if file != "" {
		c.File = file
	}
}

// Write() saves the config as JSON.
func (c *Config) Write() error {
	renameLegacy()
	c.ExercismDirectory = ""
	c.Hostname = ""
	c.ProblemsHost = ""

	// truncates existing file if it exists
	f, err := os.Create(c.File)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(c)
}

func (c *Config) configure() (*Config, error) {
	c.sanitize()

	if c.Hostname != "" {
		c.API = c.Hostname
	}

	if c.API == "" {
		c.API = hostAPI
	}

	if c.ProblemsHost != "" {
		c.XAPI = c.ProblemsHost
	}

	if c.XAPI == "" {
		c.XAPI = hostXAPI
	}

	homeDir, err := c.homeDir()
	if err != nil {
		return c, err
	}

	// use legacy value, if it exists
	if c.ExercismDirectory != "" {
		c.Dir = c.ExercismDirectory
	}

	if c.Dir == "" {
		// fall back to default value
		c.Dir = filepath.Join(homeDir, DirExercises)
	} else {
		// replace '~' with user's home
		c.Dir = strings.Replace(c.Dir, "~/", fmt.Sprintf("%s/", homeDir), 1)
	}

	if c.File == "" {
		c.File = filepath.Join(homeDir, File)
	}

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
	return filepath.Join(dir, File), nil
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

func (c *Config) sanitize() {
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.Dir = strings.TrimSpace(c.Dir)
	c.API = strings.TrimSpace(c.API)
	c.XAPI = strings.TrimSpace(c.XAPI)
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
