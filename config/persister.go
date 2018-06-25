package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ConfigType is one of the available exercism config types.
type ConfigType int

const (
	// TypeUser is configuration relevant to the individual's local setup.
	TypeUser ConfigType = iota
	// TypeCLI is configuration relevant to the API.
	TypeCLI
)

// FilePersister handles viper configs on the file system.
type FilePersister struct {
	Dir string
}

func ReadViperConfig(ct ConfigType) (*viper.Viper, error) {
	p := NewFilePersister()
	return p.Load(ct)
}

// NewFilePersister returns a persister configured with the default config directory.
func NewFilePersister() FilePersister {
	return FilePersister{
		Dir: Dir(),
	}
}

// Save stores the viper config to the configured location on the filesystem.
func (fp FilePersister) Save(v *viper.Viper, ct ConfigType) error {
	v.SetConfigType("json")
	v.AddConfigPath(fp.Dir)
	v.SetConfigName(basename(ct))

	if _, err := os.Stat(fp.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(fp.Dir, os.FileMode(0755)); err != nil {
			return err
		}
	}
	// WriteConfig is broken.
	// Someone proposed a fix in https://github.com/spf13/viper/pull/503,
	// but it doesn't work yet.
	// When it's fixed and merge we can get rid of `fp.path()`
	// and use v.WriteConfig() directly.
	return v.WriteConfigAs(fp.path(ct))
}

// Load reads a config file on the filesystem into a new Viper value.
func (fp FilePersister) Load(ct ConfigType) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(fp.Dir)
	v.SetConfigName(basename(ct))
	err := v.ReadInConfig()
	return v, err
}

func basename(ct ConfigType) string {
	switch ct {
	case TypeUser:
		return "user"
	case TypeCLI:
		return "cli"
	}
	return "unknown"
}

func (fp FilePersister) path(ct ConfigType) string {
	return filepath.Join(fp.Dir, fmt.Sprintf("%s.json", basename(ct)))
}
