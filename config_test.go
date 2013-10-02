package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
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

	ConfigToFile(tmpDir, writtenConfig)

	loadedConfig, err := ConfigFromFile(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, writtenConfig, loadedConfig)
}

func TestDemoDir(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	os.Chdir(path)

	path, err = filepath.EvalSymlinks(path)
	assert.NoError(t, err)

	path = filepath.Join(path, "exercism-demo")

	demoDir, err := DemoDirectory()
	assert.NoError(t, err)
	assert.Equal(t, demoDir, path)
}
