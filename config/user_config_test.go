// +build !windows

package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUserConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "user-config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := &UserConfig{
		Config: New(dir, "user"),
	}
	cfg.Token = "a"
	cfg.Workspace = "/a"
	cfg.APIBaseURL = "http://example.com"

	// write it
	err = cfg.Write()
	assert.NoError(t, err)

	// reload it
	cfg = &UserConfig{
		Config: New(dir, "user"),
	}
	err = cfg.Load(viper.New())
	assert.NoError(t, err)
	assert.Equal(t, "a", cfg.Token)
	assert.Equal(t, "/a", cfg.Workspace)
	assert.Equal(t, "http://example.com", cfg.APIBaseURL)
}
