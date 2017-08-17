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
		Endpoints: map[string]string{
			"a": "/a",
			"b": "/b",
		},
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
	assert.Equal(t, "/a", cfg.Endpoints["a"])
	assert.Equal(t, "/b", cfg.Endpoints["b"])
}

func TestAPIConfigSetDefaults(t *testing.T) {
	// All defaults.
	cfg := &APIConfig{}
	cfg.SetDefaults()
	assert.Equal(t, "https://api.exercism.io/v1", cfg.BaseURL)
	assert.Equal(t, "/solutions/%s", cfg.Endpoints["download"])
	assert.Equal(t, "/solutions/%s", cfg.Endpoints["submit"])

	// Override just the base url.
	cfg = &APIConfig{
		BaseURL: "http://example.com/v1",
	}
	cfg.SetDefaults()
	assert.Equal(t, "http://example.com/v1", cfg.BaseURL)
	assert.Equal(t, "/solutions/%s", cfg.Endpoints["download"])
	assert.Equal(t, "/solutions/%s", cfg.Endpoints["submit"])

	// Override just one of the endpoints.
	cfg = &APIConfig{
		Endpoints: map[string]string{
			"download": "/download/%d",
		},
	}
	cfg.SetDefaults()
	assert.Equal(t, "https://api.exercism.io/v1", cfg.BaseURL)
	assert.Equal(t, "/download/%d", cfg.Endpoints["download"])
	assert.Equal(t, "/solutions/%s", cfg.Endpoints["submit"])
}

func TestAPIConfigURL(t *testing.T) {
	cfg := &APIConfig{
		Endpoints: map[string]string{
			"a": "a/%s/a",
			"b": "b/%s/%d",
			"c": "c/%s/%s/%s",
		},
	}
	assert.Equal(t, "a/apple/a", cfg.URL("a", "apple"))
	assert.Equal(t, "b/banana/2", cfg.URL("b", "banana", 2))
	assert.Equal(t, "c/cherry/coca/cola", cfg.URL("c", "cherry", "coca", "cola"))
}
