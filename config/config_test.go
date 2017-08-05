package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type fakeConfig struct {
	*Config
	Letter string
	Number int
}

func (cfg *fakeConfig) write() error {
	return Write(cfg)
}

func (cfg *fakeConfig) load(v *viper.Viper) error {
	cfg.readIn(v)
	return v.Unmarshal(&cfg)
}

func TestFakeConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fake-config")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set the config directory to a directory that doesn't exist.
	dir := filepath.Join(tmpDir, "exercism")

	// It has access to the embedded fields.
	cfg := &fakeConfig{
		Config: New(dir, "fake"),
	}
	assert.Equal(t, dir, cfg.dir)
	assert.Equal(t, "fake", cfg.name)

	// Write the file.
	cfg.Letter = "a"
	cfg.Number = 1
	err = cfg.write()
	assert.NoError(t, err)

	// reload it
	cfg = &fakeConfig{
		Config: New(dir, "fake"),
	}
	err = cfg.load(viper.New())
	assert.NoError(t, err)

	// It wrote the non-embedded fields.
	assert.Equal(t, "a", cfg.Letter)
	assert.Equal(t, 1, cfg.Number)

	// Update the config.
	cfg.Letter = "b"
	err = cfg.write()
	assert.NoError(t, err)

	// Updating the config selectively overwrote the values.
	cfg = &fakeConfig{
		Config: New(dir, "fake"),
	}
	err = cfg.load(viper.New())
	assert.NoError(t, err)
	assert.Equal(t, "b", cfg.Letter)
	assert.Equal(t, 1, cfg.Number)
}
