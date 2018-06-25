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

func TestSetDefaultWorkspace(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	cfg := &UserConfig{Home: "/home/alice"}
	testCases := []struct {
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

	for _, tc := range testCases {
		testName := "'" + tc.in + "' should be normalized as '" + tc.out + "'"
		t.Run(testName, func(t *testing.T) {
			cfg.Workspace = tc.in
			cfg.SetDefaults()
			assert.Equal(t, tc.out, cfg.Workspace, testName)
		})
	}
}
