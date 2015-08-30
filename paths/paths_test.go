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

func TestExercises(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)
	Home = "/test/home"
	Recalculate()

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
	Recalculate()

	testCases := []struct {
		desc         string
		givenPath    string
		expectedPath string
	}{
		{
			"blank path",
			"",
			filepath.Join(Home, ".exercism.json"),
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
