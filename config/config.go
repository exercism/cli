package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/robphoenix/cli/paths"
)

const (
	// hostAPI is the endpoint to submit solutions to, and to get personalized data
	hostAPI = "http://exercism.io"
	// hostXAPI is the endpoint to fetch exercises from
	hostXAPI = "http://x.exercism.io"
)

// Config represents the settings for particular user.
// This defines both the auth for talking to the API, as well as
// where to put exercises that get downloaded.
type Config struct {
	APIKey string `json:"apiKey"`
	Dir    string `json:"dir"`
	API    string `json:"api"`
	XAPI   string `json:"xapi"`
	File   string `json:"-"` // full path to config file
}

// New returns a configuration struct with content from the exercism.json file
func New(path string) (*Config, error) {
	configPath := paths.Config(path)
	_, err := os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		if path == "" {
			configPath = paths.DefaultConfig
		}
	} else if err != nil {
		return nil, err
	}

	c := &Config{
		File: configPath,
	}
	err = c.load()
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

	if dir != "" {
		c.Dir = paths.Exercises(dir)
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

func (c *Config) load() error {
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

func (c *Config) setDefaults() error {
	if c.API == "" {
		c.API = hostAPI
	}

	if c.XAPI == "" {
		c.XAPI = hostXAPI
	}

	if _, err := url.Parse(c.API); err != nil {
		return fmt.Errorf("invalid API URL %s", err)
	}

	if _, err := url.Parse(c.XAPI); err != nil {
		return fmt.Errorf("invalid xAPI URL %s", err)
	}

	c.Dir = paths.Exercises(c.Dir)

	return nil
}
