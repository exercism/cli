package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSavingAssignment(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	assignment := Assignment{
		Track: "ruby",
		Slug:  "bob",
		Files: map[string]string{
			"bob_test.rb":      "Tests text",
			"README.md":        "Readme text",
			"/path/to/file.rb": "File text",
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
}
