package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Persister saves viper configs.
type Persister interface {
	Save(*viper.Viper, string) error
}

// FilePersister saves viper configs to the file system.
type FilePersister struct {
	Dir string
}

// Save writes the viper config to the target location on the filesystem.
func (p FilePersister) Save(v *viper.Viper, basename string) error {
	v.SetConfigType("json")
	v.AddConfigPath(p.Dir)
	v.SetConfigName(basename)

	if _, err := os.Stat(p.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(p.Dir, os.FileMode(0755)); err != nil {
			return err
		}
	}

	// WriteConfig is broken.
	// Someone proposed a fix in https://github.com/spf13/viper/pull/503,
	// but the fix doesn't work yet.
	// When it's fixed and merged we can get rid of `path`
	// and use viperConfig.WriteConfig() directly.
	path := filepath.Join(p.Dir, fmt.Sprintf("%s.json", basename))
	return v.WriteConfigAs(path)
}

// InMemoryPersister is a noop persister for use in unit tests.
type InMemoryPersister struct{}

// Save does nothing.
func (p InMemoryPersister) Save(*viper.Viper, string) error {
	return nil
}
