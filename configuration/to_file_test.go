package configuration

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	ToFile(tmpDir, writtenConfig)

	loadedConfig, err := FromFile(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, writtenConfig, loadedConfig)
}
