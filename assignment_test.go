package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSavingAssignment(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	prepareFixture(t, fmt.Sprintf("%s/ruby/bob/stub.rb", tmpDir), "Existing stub")

	assignment := Assignment{
		Track: "ruby",
		Slug:  "bob",
		Files: map[string]string{
			"bob_test.rb":     "Tests text",
			"README.md":       "Readme text",
			"path/to/file.rb": "File text",
			"stub.rb":         "New version of stub",
		},
	}

	err = SaveAssignment(tmpDir, assignment)
	assert.NoError(t, err)

	readme, err := ioutil.ReadFile(tmpDir + "/ruby/bob/README.md")
	assert.NoError(t, err)
	assert.Equal(t, string(readme), "Readme text")

	tests, err := ioutil.ReadFile(tmpDir + "/ruby/bob/bob_test.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(tests), "Tests text")

	fileInDir, err := ioutil.ReadFile(tmpDir + "/ruby/bob/path/to/file.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(fileInDir), "File text")

	stubFile, err := ioutil.ReadFile(tmpDir + "/ruby/bob/stub.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(stubFile), "Existing stub")
}

func prepareFixture(t *testing.T, fixture, s string) {
	err := os.MkdirAll(filepath.Dir(fixture), 0755)
	assert.NoError(t, err)

	err = ioutil.WriteFile(fixture, []byte(s), 0644)
	assert.NoError(t, err)

	// ensure fixture is set up correctly
	fixtureContents, err := ioutil.ReadFile(fixture)
	assert.NoError(t, err)
	assert.Equal(t, string(fixtureContents), s)
}
