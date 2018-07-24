package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
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

	// We're going to load up a viper that is bound to some command-line flags.
	// First we need flags.
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flagSet.IntP("number", "n", 0, "pick a number, any number")
	flagSet.StringP("letter", "l", "", "something from a nice alphabet")
	flagSet.Set("number", "1")
	flagSet.Set("letter", "a")

	// Bind the flags to a new viper value.
	v := viper.New()
	v.BindPFlag("number", flagSet.Lookup("number"))
	v.BindPFlag("letter", flagSet.Lookup("letter"))

	// Binding the flags loaded the values into viper.
	assert.Equal(t, 1, v.Get("number"))
	assert.Equal(t, "a", v.Get("letter"))

	// Load viper into the config value.
	err = cfg.load(v)
	assert.NoError(t, err)

	// The original flag values have been loaded into the struct value.
	assert.Equal(t, 1, cfg.Number)
	assert.Equal(t, "a", cfg.Letter)

	// Write the file.
	err = cfg.write()
	assert.NoError(t, err)

	// Reload it.
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

func TestInferSiteURL(t *testing.T) {
	testCases := []struct {
		api, url string
	}{
		{"https://api.exercism.io/v1", "https://exercism.io"},
		{"https://v2.exercism.io/api/v1", "https://v2.exercism.io"},
		{"https://mentors-beta.exercism.io/api/v1", "https://mentors-beta.exercism.io"},
		{"http://localhost:3000/api/v1", "http://localhost:3000"},
		{"", "https://exercism.io"},            // use the default
		{"http://whatever", "http://whatever"}, // you're on your own, pal
	}

	for _, tc := range testCases {
		assert.Equal(t, InferSiteURL(tc.api), tc.url)
	}
}
