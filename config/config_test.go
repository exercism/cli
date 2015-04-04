package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.Link(fixturePath(t, "config.json"), configPath); err != nil {
		t.Fatal(err)
	}
	dirtyPath := filepath.Join(tmpDir, "dirty.json")
	if err := os.Link(fixturePath(t, "dirty.json"), dirtyPath); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		desc                string
		in                  string // the name of the file passed as a command line argument
		out                 string // the name of the file that the config will be written to
		dir, key, api, xapi string // the actual config values
	}{
		{
			desc: "defaults",
			in:   "",
			out:  filepath.Join(tmpDir, File),
			dir:  filepath.Join(tmpDir, DirExercises),
			key:  "",
			api:  hostAPI,
			xapi: hostXAPI,
		},
		{
			desc: "no such file",
			in:   filepath.Join(tmpDir, "no-such.json"),
			out:  filepath.Join(tmpDir, "no-such.json"),
			dir:  filepath.Join(tmpDir, DirExercises),
			key:  "",
			api:  hostAPI,
			xapi: hostXAPI,
		},
		{
			desc: "file exists",
			in:   configPath,
			out:  configPath,
			dir:  "/a/b/c",
			key:  "abc123",
			api:  "http://api.example.com",
			xapi: "http://x.example.com",
		},
		{
			desc: "unexpanded path",
			in:   "~/config.json",
			out:  configPath,
			dir:  "/a/b/c",
			key:  "abc123",
			api:  "http://api.example.com",
			xapi: "http://x.example.com",
		},
		{
			desc: "sanitizes whitespace",
			in:   "~/dirty.json",
			out:  filepath.Join(tmpDir, "dirty.json"),
			dir:  "/a/b/c",
			key:  "abc123",
			api:  "http://api.example.com",
			xapi: "http://x.example.com",
		},
	}

	for _, tc := range testCases {
		c := &Config{home: tmpDir}

		if err := c.load(tc.in); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.out, c.File, tc.desc)
		assert.Equal(t, tc.dir, c.Dir, tc.desc)
		assert.Equal(t, tc.key, c.APIKey, tc.desc)
		assert.Equal(t, tc.api, c.API, tc.desc)
		assert.Equal(t, tc.xapi, c.XAPI, tc.desc)
	}
}

func TestReadDirectory(t *testing.T) {
	// if the provided path is a directory, append the default filename
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	myConfig, err := New(tmpDir)
	assert.NoError(t, err)

	expected := filepath.Join(tmpDir, File)
	actual := myConfig.File
	assert.Equal(t, expected, actual)

	// if it can't determine if the provided path is a directory, don't modify
	// the path
	myConfig, err = New("badpath")
	assert.NoError(t, err)
	assert.Equal(t, "badpath", myConfig.File)
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	invalidPath := filepath.Join(tmpDir, "config_invalid.json")
	if err := os.Link(fixturePath(t, "config_invalid.json"), invalidPath); err != nil {
		t.Fatal(err)
	}
	c := &Config{home: tmpDir}

	err = c.load("~/config_invalid.json")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid JSON syntax")
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
		XAPI:   "localhost",
		File:   filename,
	}

	c1.Write()

	c2, err := New(filename)
	assert.NoError(t, err)

	assert.Equal(t, c1.APIKey, c2.APIKey)
	assert.Equal(t, c1.Dir, c2.Dir)
	assert.Equal(t, c1.API, c2.API)
	assert.Equal(t, c1.XAPI, c2.XAPI)
}

func TestUpdateConfig(t *testing.T) {
	c := &Config{
		APIKey: "MyKey",
		API:    "localhost",
		Dir:    "/exercism/directory",
		XAPI:   "localhost",
	}

	// Test the blank values don't overwrite existing values
	c.Update("", "", "", "")
	assert.Equal(t, "MyKey", c.APIKey)
	assert.Equal(t, "localhost", c.API)
	assert.Equal(t, "/exercism/directory", c.Dir)
	assert.Equal(t, "localhost", c.XAPI)

	// Test that each value can be overwritten
	c.Update("NewKey", "http://example.com", "/tmp/exercism", "http://x.example.org")
	assert.Equal(t, "NewKey", c.APIKey)
	assert.Equal(t, "http://example.com", c.API)
	assert.Equal(t, "/tmp/exercism", c.Dir)
	assert.Equal(t, "http://x.example.org", c.XAPI)
}

func fixturePath(t *testing.T, filename string) string {
	_, caller, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	return filepath.Join(filepath.Dir(caller), "..", "fixtures", filename)
}
