package exercism

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func assertFileExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err != nil {
		t.Errorf("File [%s] does not exist.", filename)
	}
}

func asserFileDoesNotExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err == nil {
		t.Errorf("File [%s] already exist.", filename)
	}
}

func TestLoginCreatesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	filename := tmpDir + "/" + FILENAME
	asserFileDoesNotExist(t, filename)

	config := Config{"githubUser", "MyApiKey", "/my/exercism/directory"}

	Login(tmpDir, config)

	assertFileExists(t, filename)

	contents, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)

	configEntries := &Config{}
	err = json.Unmarshal(contents, configEntries)
	assert.NoError(t, err)

	assert.Equal(t, configEntries.GithubUsername, "githubUser")
	assert.Equal(t, configEntries.ApiKey, "MyApiKey")
	assert.Equal(t, configEntries.ExercismDirectory, "/my/exercism/directory")
}

func TestLogoutDeletesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	filename := tmpDir + "/" + FILENAME
	err = ioutil.WriteFile(filename, []byte("exercism config\n"), 0644)
	assert.NoError(t, err)

	Logout(tmpDir)

	asserFileDoesNotExist(t, filename)
}
