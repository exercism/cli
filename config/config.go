package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

// New returns a configuration struct with content from the exercism.json file
func New(path string) (*Config, error) {
	c := &Config{}
	err := c.load(path)
	return c, err
}

// Update sets new values where given.
func (c *Config) Update(key, host, dir, xapi string) error {
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
		if err := c.SetDir(dir); err != nil {
			return err
		}
	}

	xapi = strings.TrimSpace(xapi)
	if xapi != "" {
		c.XAPI = xapi
	}

	return nil
}

// Write saves the config as JSON.
func (c *Config) Write() error {
	// truncates existing file if it exists
	f, err := os.Create(c.File)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}

func (c *Config) load(argPath string) error {
	path, err := c.resolvePath(argPath)
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

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		var extra string
		if serr, ok := err.(*json.SyntaxError); ok {
			if _, serr := f.Seek(0, os.SEEK_SET); serr != nil {
				log.Fatalf("seek error: %v", serr)
			}
			line, str := findInvalidJSON(f, serr.Offset)
			extra = fmt.Sprintf(":\ninvalid JSON syntax at line %d:\n%s",
				line, str)
		}
		return fmt.Errorf("error parsing JSON in the config file %s%s\n%s", f.Name(), extra, err)
	}

	return nil
}

func findInvalidJSON(f io.Reader, pos int64) (int, string) {
	var (
		col     int
		line    int
		errLine []byte
	)
	buf := new(bytes.Buffer)
	fb := bufio.NewReader(f)

	for c := int64(0); c < pos; {
		b, err := fb.ReadBytes('\n')
		if err != nil {
			log.Fatalf("read error: %v", err)
		}
		c += int64(len(b))
		col = len(b) - int(c-pos)

		line++
		errLine = b
	}

	if len(errLine) != 0 {
		buf.WriteString(fmt.Sprintf("%5d: %s <~", line, errLine[:col]))
	}

	return line, buf.String()
}

// IsAuthenticated returns true if the config contains an API key.
// This does not check whether or not that key is valid.
func (c *Config) IsAuthenticated() bool {
	return c.APIKey != ""
}

// homeDir caches the lookup of the user's home directory.
func (c *Config) homeDir() (string, error) {
	if c.home != "" {
		return c.home, nil // only set during testing
	}
	return Home()
}

func (c *Config) resolvePath(argPath string) (string, error) {
	path := argPath
	if path == "" {
		path = filepath.Join("~", File)
	}
	h, err := c.homeDir()
	if err != nil {
		return "", err
	}
	path = expandHome(path, h)

	fi, _ := os.Stat(path)
	if fi != nil && fi.IsDir() {
		path = filepath.Join(path, File)
	}

	return path, nil
}

func (c *Config) setDefaults() error {
	if c.API == "" {
		c.API = hostAPI
	}

	if c.XAPI == "" {
		c.XAPI = hostXAPI
	}

	if err := c.SetDir(c.Dir); err != nil {
		return err
	}

	return nil
}

// SetDir sets the configuration directory to the given path
// or defaults to the home exercism directory
func (c *Config) SetDir(path string) error {
	home, err := c.homeDir()
	if err != nil {
		return err
	}

	var dir string

	if path == "" {
		dir = filepath.Join(home, DirExercises)
	} else {
		dir = path
	}

	dir = expandHome(dir, home)

	// if the user has provided us with a relative path, make it absolute so
	// it will always work
	if !filepath.IsAbs(dir) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = filepath.Join(wd, dir)
	}

	c.Dir = dir

	return nil
}

func expandHome(path, home string) string {
	if path[:2] == "~"+string(os.PathSeparator) {
		return strings.Replace(path, "~", home, 1)
	}
	return path
}
