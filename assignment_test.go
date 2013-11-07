package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSavingAssignment(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	assignment := Assignment{
		Track:    "ruby",
		Slug:     "bob",
		Readme:   "Readme text",
		StubFile: "bob.rb",
		Stub:     "Stub Text",
		TestFile: "bob_test.rb",
		Tests:    "Tests Text",
	}

	err = SaveAssignment(tmpDir, assignment)
	assert.NoError(t, err)

	readme, err := ioutil.ReadFile(tmpDir + "/ruby/bob/README.md")
	assert.NoError(t, err)
	assert.Equal(t, string(readme), "Readme text")

	stub, err := ioutil.ReadFile(tmpDir + "/ruby/bob/bob.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(stub), "Stub Text")

	tests, err := ioutil.ReadFile(tmpDir + "/ruby/bob/bob_test.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(tests), "Tests Text")

	assignment = Assignment{
		Track:    "ruby",
		Slug:     "space-age",
		Readme:   "Readme text",
		StubFile: "",
		Stub:     "",
		TestFile: "space-age_test.rb",
		Tests:    "Tests Text",
	}

	_, err = ioutil.ReadFile(tmpDir + "/ruby/space-age/space-age.rb")
	assert.True(t, os.IsNotExist(err))
}
