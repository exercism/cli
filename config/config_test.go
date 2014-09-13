package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValues(t *testing.T) {
	c := &Config{}
	c.home = "/home/alice"
	c.configure()
	assert.Equal(t, "", c.APIKey)
	assert.Equal(t, "http://exercism.io", c.Hostname)
	assert.Equal(t, "/home/alice/exercism", c.Dir)
}

func TestCustomValues(t *testing.T) {
	c := &Config{
		APIKey:   "abc123",
		Hostname: "http://example.org",
		Dir:      "/path/to/exercises",
	}
	c.configure()
	assert.Equal(t, "abc123", c.APIKey)
	assert.Equal(t, "http://example.org", c.Hostname)
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
		APIKey:   "   abc123\n\r\n  ",
		Hostname: "       ",
		Dir:      "  \r\n/path/to/exercises   \r\n",
	}
	c.configure()
	assert.Equal(t, "abc123", c.APIKey)
	assert.Equal(t, "http://exercism.io", c.Hostname)
	assert.Equal(t, "/path/to/exercises", c.Dir)
}

func TestFilePath(t *testing.T) {
	// defaults to home directory
	c := &Config{}
	c.home = "/home/alice"
	c.configure()
	assert.Equal(t, "/home/alice/.exercism.json", c.File())

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
	assert.Equal(t, c.Hostname, "http://exercism.io")
	assert.False(t, c.IsAuthenticated())
	if !strings.HasSuffix(c.Dir, "/exercism") {
		t.Fatal("Default unconfigured config should use home dir")
	}
}

func TestReadingWritingConfig(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	filename := fmt.Sprintf("%s/%s", tmpDir, File)
	assert.NoError(t, err)

	c1 := &Config{
		APIKey:   "MyKey",
		Dir:      "/exercism/directory",
		Hostname: "localhost",
	}

	c1.SavePath(filename)
	c1.Write()

	c2, err := Read(filename)
	assert.NoError(t, err)

	assert.Equal(t, c1.APIKey, c2.APIKey)
	assert.Equal(t, c1.Dir, c2.Dir)
	assert.Equal(t, c1.Hostname, c2.Hostname)
}

func TestDecodingConfig(t *testing.T) {
	unsanitizedJSON := `{"apiKey":"MyKey  ","exercismDirectory":"/exercism/directory\r\n","hostname":"localhost \r\n"}`
	c1 := &Config{
		APIKey:   "MyKey",
		Dir:      "/exercism/directory",
		Hostname: "localhost",
	}
	b := bytes.NewBufferString(unsanitizedJSON)
	c2, err := Decode(b)
	assert.NoError(t, err)

	assert.Equal(t, c1.APIKey, c2.APIKey)
	assert.Equal(t, c1.Dir, c2.Dir)
	assert.Equal(t, c1.Hostname, c2.Hostname)
}

func TestEncodingConfig(t *testing.T) {
	c := Config{
		APIKey:       "MyKey ",
		Dir:          "/home/user name  ",
		Hostname:     "localhost  ",
		ProblemsHost: " localhost:9292 ",
	}
	c.configure()
	sanitizedJSON := `{"apiKey":"MyKey","exercismDirectory":"/home/user name","hostname":"localhost","problemsHost":"localhost:9292"}
`

	buf := new(bytes.Buffer)
	err := c.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, sanitizedJSON, buf.String())
}
