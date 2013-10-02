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
		Track:    "ruby",
		Slug:     "bob",
		Readme:   "Readme text",
		TestFile: "bob_test.rb",
		Tests:    "Tests Text",
	}

	err = SaveAssignment(tmpDir, assignment)
	assert.NoError(t, err)

	readme, err := ioutil.ReadFile(tmpDir + "/ruby/bob/README.md")
	assert.NoError(t, err)
	assert.Equal(t, string(readme), "Readme text")

	tests, err := ioutil.ReadFile(tmpDir + "/ruby/bob/bob_test.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(tests), "Tests Text")
}
