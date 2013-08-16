package exercism

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/user"
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

	c := Config{}

	user := user.User{HomeDir: "/Users/foo"}
	ConfigToFile(user, tmpDir, c)

	Logout(tmpDir)

	asserFileDoesNotExist(t, configFilename(tmpDir))
}
