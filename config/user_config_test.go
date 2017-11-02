// +build !windows

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
}

func TestNormalizeWorkspace(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	cfg := &UserConfig{Home: "/home/alice"}
	tests := []struct {
		in, out string
	}{
		{"", ""}, // don't make wild guesses
		{"/home/alice///foobar", "/home/alice/foobar"},
		{"~/foobar", "/home/alice/foobar"},
		{"/foobar/~/noexpand", "/foobar/~/noexpand"},
		{"/no/modification", "/no/modification"},
		{"relative", filepath.Join(cwd, "relative")},
		{"relative///path", filepath.Join(cwd, "relative", "path")},
	}

	for _, test := range tests {
		cfg.Workspace = test.in
		cfg.Normalize()
		assert.Equal(t, test.out, cfg.Workspace)
	}
}
