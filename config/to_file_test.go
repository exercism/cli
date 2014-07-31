package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadingWritingConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	filename := Filename(tmpDir)
	assert.NoError(t, err)

	writtenConfig := Config{
		GithubUsername:    "user",
		ApiKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
	}

	ToFile(filename, writtenConfig)

	loadedConfig, err := FromFile(filename)
	assert.NoError(t, err)

	assert.Equal(t, writtenConfig, loadedConfig)
}
