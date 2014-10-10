package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValues(t *testing.T) {
	c := &Config{}
	c.home = "/home/alice"
	c.configure()
	assert.Equal(t, "", c.APIKey)
	assert.Equal(t, "http://exercism.io", c.API)
	assert.Equal(t, filepath.FromSlash("/home/alice/exercism"), c.Dir)
}

func TestCustomValues(t *testing.T) {
	c := &Config{
		APIKey: "abc123",
		API:    "http://example.org",
		Dir:    "/path/to/exercises",
	}
	c.configure()
	assert.Equal(t, "abc123", c.APIKey)
	assert.Equal(t, "http://example.org", c.API)
	assert.Equal(t, "/path/to/exercises", c.Dir)
}

func TestExpandHomeDir(t *testing.T) {
	c := &Config{Dir: "~/practice"}
	c.home = "/home/alice"
	c.configure()
	assert.Equal(t, "/home/alice/practice", c.Dir)
}

func TestSanitizeWhitespace(t *testing.T) {
	c := &Config{
		APIKey: "   abc123\n\r\n  ",
		API:    "       ",
		Dir:    "  \r\n/path/to/exercises   \r\n",
	}
	c.configure()
	assert.Equal(t, "abc123", c.APIKey)
	assert.Equal(t, "http://exercism.io", c.API)
	assert.Equal(t, "/path/to/exercises", c.Dir)
}

func TestFilePath(t *testing.T) {
	// defaults to home directory
	c := &Config{}
	c.home = "/home/alice"
	c.configure()
	assert.Equal(t, filepath.FromSlash("/home/alice/.exercism.json"), c.File())

	// can override location of config file
	c = &Config{}
	c.configure()
	c.SavePath("/tmp/config/exercism.conf")
	assert.Equal(t, "/tmp/config/exercism.conf", c.File())
}

func TestReadNonexistantConfig(t *testing.T) {
	c, err := Read("/no/such/config.json")
	assert.NoError(t, err)
	assert.Equal(t, c.APIKey, "")
	assert.Equal(t, c.API, "http://exercism.io")
	assert.False(t, c.IsAuthenticated())
	if !strings.HasSuffix(c.Dir, filepath.FromSlash("/exercism")) {
		t.Fatal("Default unconfigured config should use home dir")
	}
}

func TestReadingWritingConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	filename := fmt.Sprintf("%s/%s", tmpDir, File)
	assert.NoError(t, err)

	c1 := &Config{
		APIKey: "MyKey",
		Dir:    "/exercism/directory",
		API:    "localhost",
	}
	c1.configure()

	c1.SavePath(filename)
	c1.Write()

	c2, err := Read(filename)
	assert.NoError(t, err)

	assert.Equal(t, c1.APIKey, c2.APIKey)
	assert.Equal(t, c1.Dir, c2.Dir)
	assert.Equal(t, c1.API, c2.API)
}

func TestUpdateConfig(t *testing.T) {
	c := &Config{
		APIKey: "MyKey",
		Dir:    "/exercism/directory",
		API:    "localhost",
	}

	c.Update("NewKey", "", "")
	assert.Equal(t, "NewKey", c.APIKey)
	assert.Equal(t, "localhost", c.API)
	assert.Equal(t, "/exercism/directory", c.Dir)

	c.Update("", "http://example.com", "")
	assert.Equal(t, "NewKey", c.APIKey)
	assert.Equal(t, "http://example.com", c.API)
	assert.Equal(t, "/exercism/directory", c.Dir)

	c.Update("", "", "/tmp/exercism")
	assert.Equal(t, "NewKey", c.APIKey)
	assert.Equal(t, "http://example.com", c.API)
	assert.Equal(t, "/tmp/exercism", c.Dir)
}

func TestReadDefaultConfig(t *testing.T) {
	dir, err := filepath.Abs("../fixtures/home")
	assert.NoError(t, err)

	c := &Config{home: dir}
	err = c.Read("")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", c.APIKey)
	assert.Equal(t, "/path/to/exercism", c.Dir)
	assert.Equal(t, "http://example.com", c.API)
	assert.Equal(t, "http://x.example.com", c.XAPI)
}

func TestReadCustomConfig(t *testing.T) {
	dir, err := filepath.Abs("../fixtures/home/")
	assert.NoError(t, err)

	c := &Config{home: dir}
	file := fmt.Sprintf("%s/custom.json", dir)
	err = c.Read(file)
	assert.NoError(t, err)
	assert.Equal(t, "xyz000", c.APIKey)
	assert.Equal(t, "/tmp/exercism", c.Dir)
	assert.Equal(t, "http://example.org", c.API)
	assert.Equal(t, "http://x.example.org", c.XAPI)
}

func TestReadLegacyConfig(t *testing.T) {
	dir, err := filepath.Abs("../fixtures/home/")
	assert.NoError(t, err)

	c := &Config{home: dir}
	file := fmt.Sprintf("%s/legacy.json", dir)
	err = c.Read(file)
	assert.NoError(t, err)
	assert.Equal(t, "prq567", c.APIKey)
	assert.Equal(t, "/tmp/stuff", c.Dir)
	assert.Equal(t, "http://api.example.com", c.API)
	assert.Equal(t, "http://problems.example.com", c.XAPI)
}
