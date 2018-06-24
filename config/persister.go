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
	// TypeAPI is configuration relevant to the API.
	TypeAPI
)

// Persister saves viper configurations.
type Persister interface {
	Load(ConfigType) (*viper.Viper, error)
	Save(*viper.Viper, ConfigType) error
}

// FilePersister handles viper configs on the file system.
type FilePersister struct {
	Dir string
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

// InMemoryPersister stores viper configs in memory (for testing).
type InMemoryPersister struct {
	store map[string]*viper.Viper
}

// NewInMemoryPersister creates a new in memory store.
func NewInMemoryPersister() *InMemoryPersister {
	return &InMemoryPersister{
		store: map[string]*viper.Viper{},
	}
}

// Save stores the viper config to memory.
func (imp *InMemoryPersister) Save(v *viper.Viper, ct ConfigType) error {
	imp.store[basename(ct)] = v
	return nil
}

// Load returns a viper config stored in memory.
func (imp *InMemoryPersister) Load(ct ConfigType) (*viper.Viper, error) {
	return imp.store[basename(ct)], nil
}

func basename(ct ConfigType) string {
	switch ct {
	case TypeUser:
		return "user"
	case TypeAPI:
		return "api"
	}
	return "unknown"
}

func (fp FilePersister) path(ct ConfigType) string {
	return filepath.Join(fp.Dir, fmt.Sprintf("%s.json", basename(ct)))
}
