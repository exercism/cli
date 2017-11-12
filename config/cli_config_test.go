package config

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCLIConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "cli-config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := &CLIConfig{
		Config: New(dir, "cli"),
		Tracks: Tracks{
			"bogus": &Track{ID: "bogus"},
			"fake":  &Track{ID: "fake", IgnorePatterns: []string{"c", "b", "a"}},
		},
	}

	// write it
	err = cfg.Write()
	assert.NoError(t, err)

	// reload it
	cfg = &CLIConfig{
		Config: New(dir, "cli"),
	}
	err = cfg.Load(viper.New())
	assert.NoError(t, err)
	assert.Equal(t, "bogus", cfg.Tracks["bogus"].ID)
	assert.Equal(t, "fake", cfg.Tracks["fake"].ID)

	// The ignore patterns got sorted.
	expected := append(defaultIgnorePatterns, "a", "b", "c")
	sort.Strings(expected)
	assert.Equal(t, expected, cfg.Tracks["fake"].IgnorePatterns)
}

func TestCLIConfigValidate(t *testing.T) {
	cfg := &CLIConfig{
		Tracks: Tracks{
			"fake": &Track{
				ID:             "fake",
				IgnorePatterns: []string{"(?=re)"}, // not a valid regex
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)
}

func TestCLIConfigSetDefaults(t *testing.T) {
	// No tracks, no defaults.
	cfg := &CLIConfig{}
	cfg.SetDefaults()
	assert.Equal(t, &CLIConfig{}, cfg)

	// With a track, gets defaults.
	cfg = &CLIConfig{
		Tracks: map[string]*Track{
			"bogus": {
				ID: "bogus",
			},
		},
	}
	cfg.SetDefaults()
	assert.Equal(t, defaultIgnorePatterns, cfg.Tracks["bogus"].IgnorePatterns)

	// With partial defaults and extras, gets everything.
	cfg = &CLIConfig{
		Tracks: map[string]*Track{
			"bogus": {
				ID:             "bogus",
				IgnorePatterns: []string{"[.]solution[.]json", "_spec[.]ext$"},
			},
		},
	}
	cfg.SetDefaults()
	expected := append(defaultIgnorePatterns, "_spec[.]ext$")
	sort.Strings(expected)
	assert.Equal(t, expected, cfg.Tracks["bogus"].IgnorePatterns)
}
