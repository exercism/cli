package paths

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	assert.Equal(t, os.Getenv("HOME"), Home)
}

func TestConfigHome(t *testing.T) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		assert.Equal(t, filepath.Join(Home, ".config"), ConfigHome)
	} else {
		assert.Equal(t, xdgConfigHome, ConfigHome)
	}
}

func TestExercises(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)
	Home = "/test/home"

	testCases := []struct {
		givenPath    string
		expectedPath string
	}{
		{"", "/test/home/exercism"},
		{"~/foobar", "/test/home/foobar"},
		{"/foobar/~/noexpand", "/foobar/~/noexpand"},
		{"/no/modification", "/no/modification"},
		{"relativePath", filepath.Join(dir, "relativePath")},
	}

	for _, testCase := range testCases {
		actual := Exercises(testCase.givenPath)
		assert.Equal(t, testCase.expectedPath, actual)
	}
}

func TestConfig(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	Home = dir
	ConfigHome = dir

	testCases := []struct {
		desc         string
		givenPath    string
		expectedPath string
	}{
		{
			"blank path",
			"",
			filepath.Join(ConfigHome, File),
		},
		{
			"unknown path is expanded, but not modified",
			"~/unknown",
			filepath.Join(Home, "unknown"),
		},
		{
			"absolute path is unmodified",
			Config(Config("")),
			Config(""),
		},
		{
			"dir path has the config file appended",
			dir,
			filepath.Join(dir, File),
		},
	}

	for _, tc := range testCases {
		actual := Config(tc.givenPath)
		assert.Equal(t, tc.expectedPath, actual, tc.desc)
	}
}
