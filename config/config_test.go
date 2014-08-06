package config

import (
	"bytes"
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

	writtenConfig := Config{
		GithubUsername:    "user",
		APIKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
		Hostname:          "localhost",
	}

	ToFile(filename, writtenConfig)

	loadedConfig, err := FromFile(filename)
	assert.NoError(t, err)

	assert.Equal(t, writtenConfig, loadedConfig)
}

func TestDecodingConfig(t *testing.T) {
	unsanitizedJson := `{"githubUsername":"user ","apiKey":"MyKey  ","exercismDirectory":"/exercism/directory\r\n","hostname":"localhost \r\n"}`
	sanitizedConfig := Config{
		GithubUsername:    "user",
		APIKey:            "MyKey",
		ExercismDirectory: "/exercism/directory",
		Hostname:          "localhost",
	}
	b := bytes.NewBufferString(unsanitizedJson)
	c, err := Decode(b)

	assert.NoError(t, err)
	assert.Equal(t, sanitizedConfig, c)
}

func TestEncodingConfig(t *testing.T) {
	currentConfig := Config{
		GithubUsername:    "user\r\n",
		APIKey:            "MyKey ",
		ExercismDirectory: "/home/user name  ",
		Hostname:          "localhost  ",
	}
	sanitizedJson := `{"githubUsername":"user","apiKey":"MyKey","exercismDirectory":"/home/user name","hostname":"localhost"}
`

	buf := new(bytes.Buffer)
	err := Encode(buf, currentConfig)

	assert.NoError(t, err)
	assert.Equal(t, sanitizedJson, buf.String())
}
