package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
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
	return ioutil.WriteFile(f.File(), b, os.FileMode(0644))
}
