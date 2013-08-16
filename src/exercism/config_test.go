package exercism

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os/user"
	"testing"
)

func TestReadingWritingConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	writtenConfig := Config{
		GithubUsername:    "user",
		ApiKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
	}

	u := user.User{HomeDir: "/Users/foo"}
	ConfigToFile(u, tmpDir, writtenConfig)

	loadedConfig, err := ConfigFromFile(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, writtenConfig, loadedConfig)
}

func TestExpandsTildeInExercismDirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	writtenConfig := Config{
		GithubUsername:    "user",
		ApiKey:            "MyKey",
		ExercismDirectory: "~/exercism/directory",
	}

	u := user.User{HomeDir: "/Users/foo"}
	ConfigToFile(u, tmpDir, writtenConfig)

	loadedConfig, err := ConfigFromFile(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, loadedConfig.ExercismDirectory, "/Users/foo/exercism/directory")
}
