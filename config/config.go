package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
)

var (
	defaultBaseURL = "https://api.exercism.io/v1"
)

// Config is a wrapper around a viper configuration.
type Config struct {
	dir  string
	name string
}

// New creates a default config value for the given directory.
func New(dir, name string) *Config {
	return &Config{
		dir:  dir,
		name: name,
	}
}

// File is the full path to the config file.
func (cfg *Config) File() string {
	return filepath.Join(cfg.dir, fmt.Sprintf("%s.json", cfg.name))
}

func (cfg *Config) readIn(v *viper.Viper) {
	v.AddConfigPath(cfg.dir)
	v.SetConfigName(cfg.name)
	v.SetConfigType("json")
	v.ReadInConfig()
}

type filer interface {
	File() string
}

// Write stores the config into a file.
func Write(f filer) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if err := ensureDir(f); err != nil {
		return err
	}
	return ioutil.WriteFile(f.File(), b, os.FileMode(0644))
}

func ensureDir(f filer) error {
	dir := filepath.Dir(f.File())
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, os.FileMode(0755))
	}
	return err
}

// InferSiteURL guesses what the website URL is.
// The basis for the guess is which API we're submitting to.
func InferSiteURL(apiURL string) string {
	if apiURL == "" {
		apiURL = defaultBaseURL
	}
	if apiURL == "https://api.exercism.io/v1" {
		return "https://exercism.io"
	}
	re := regexp.MustCompile("^(https?://[^/]*).*")
	return re.ReplaceAllString(apiURL, "$1")
}
