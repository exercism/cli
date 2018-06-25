package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAPIConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "api-config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := &APIConfig{
		Config:  New(dir, "api"),
		BaseURL: "http://example.com/v1",
	}

	// write it
	err = cfg.Write()
	assert.NoError(t, err)

	// reload it
	cfg = &APIConfig{
		Config: New(dir, "api"),
	}
	err = cfg.Load(viper.New())
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/v1", cfg.BaseURL)
}

func TestAPIConfigSetDefaults(t *testing.T) {
	// All defaults.
	cfg := &APIConfig{}
	cfg.SetDefaults()
	assert.Equal(t, "https://v2.exercism.io/api/v1", cfg.BaseURL)

	// Override just the base url.
	cfg = &APIConfig{
		BaseURL: "http://example.com/v1",
	}
	cfg.SetDefaults()
	assert.Equal(t, "http://example.com/v1", cfg.BaseURL)
}
