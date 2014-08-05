package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDemoDir(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	os.Chdir(path)

	path, err = filepath.EvalSymlinks(path)
	assert.NoError(t, err)

	path = filepath.Join(path, "exercism-demo")

	demoDir, err := demoDirectory()
	assert.NoError(t, err)
	assert.Equal(t, demoDir, path)
}

func TestExpandsTildeInExercismDirectory(t *testing.T) {
	expandedDir := ReplaceTilde("~/exercism/directory")
	assert.NotContains(t, "~", expandedDir)
}

func TestReadingWritingConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	filename := Filename(tmpDir)
	assert.NoError(t, err)

	currentConfig := Config{
		GithubUsername:    "user",
		APIKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
		Hostname:          "localhost\r\n",
	}
	sanitizedConfig := Config{
		GithubUsername:    "user",
		APIKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
		Hostname:          "localhost",
	}

	ToFile(filename, currentConfig)

	loadedConfig, err := FromFile(filename)
	assert.NoError(t, err)

	assert.Equal(t, sanitizedConfig, loadedConfig)
}

func TestSanitizeFields(t *testing.T) {
	config := Config{
		GithubUsername:    "user ",
		APIKey:            "MyKey     ",
		ExercismDirectory: "/home/user name\r\n",
		Hostname:          "localhost\n",
	}
	sanitizedConfig := Config{
		GithubUsername:    "user",
		APIKey:            "MyKey",
		ExercismDirectory: "/home/user name",
		Hostname:          "localhost",
	}
	sanitize(&config)

	assert.Equal(t, config, sanitizedConfig)
}
