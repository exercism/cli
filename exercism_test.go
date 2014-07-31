package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func assertFileDoesNotExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err == nil {
		t.Errorf("File [%s] already exist.", filename)
	}
}

func TestLogoutDeletesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	c := config.Config{}

	config.ToFile(tmpDir, c)

	logout(tmpDir)

	assertFileDoesNotExist(t, config.Filename(tmpDir))
}

func TestAskForConfigInfoAllowsSpaces(t *testing.T) {
	oldStdin := os.Stdin
	dirName := "dirname with spaces"
	userName := "TestUsername"
	apiKey := "abc123"

	fakeStdin, err := ioutil.TempFile("", "stdin_mock")
	assert.NoError(t, err)

	fakeStdin.WriteString(fmt.Sprintf("%s\r\n%s\r\n%s\r\n", userName, apiKey, dirName))
	assert.NoError(t, err)

	_, err = fakeStdin.Seek(0, os.SEEK_SET)
	assert.NoError(t, err)

	defer fakeStdin.Close()

	os.Stdin = fakeStdin

	c, err := askForConfigInfo()
	if err != nil {
		t.Errorf("Error asking for configuration info [%v]", err)
	}
	os.Stdin = oldStdin
	absoluteDirName, _ := absolutePath(dirName)
	_, err = os.Stat(absoluteDirName)
	if err != nil {
		t.Errorf("Excercism directory [%s] was not created.", absoluteDirName)
	}
	os.Remove(absoluteDirName)
	os.Remove(fakeStdin.Name())

	assert.Equal(t, c.ExercismDirectory, absoluteDirName)
	assert.Equal(t, c.GithubUsername, userName)
	assert.Equal(t, c.APIKey, apiKey)
}
