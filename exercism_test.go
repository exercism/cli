package main

import (
	"github.com/exercism/cli/configuration"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func asserFileDoesNotExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err == nil {
		t.Errorf("File [%s] already exist.", filename)
	}
}

func TestLogoutDeletesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	c := configuration.Config{}

	configuration.ToFile(tmpDir, c)

	logout(tmpDir)

	asserFileDoesNotExist(t, configuration.Filename(tmpDir))
}
