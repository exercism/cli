package workspace

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	testConfig, ok := TestConfigurations["elixir"]
	assert.True(t, ok, "unexpectedly unable to find elixir test config")

	cmd, err := testConfig.GetTestCommand()
	assert.NoError(t, err)

	assert.Equal(t, cmd, "mix test")
}

func TestWindowsCommands(t *testing.T) {
	testConfig, ok := TestConfigurations["cobol"]
	assert.True(t, ok, "unexpectedly unable to find cobol test config")

	cmd, err := testConfig.GetTestCommand()
	assert.NoError(t, err)

	if runtime.GOOS == "windows" {
		assert.Contains(t, cmd, ".ps1")
		assert.NotContains(t, cmd, ".sh")
	} else {
		assert.Contains(t, cmd, ".sh")
		assert.NotContains(t, cmd, ".ps1")
	}
}

func TestGetCommandMissingConfig(t *testing.T) {
	testConfig, ok := TestConfigurations["ruby"]
	assert.True(t, ok, "unexpectedly unable to find ruby test config")

	_, err := testConfig.GetTestCommand()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), ".exercism/config.json: no such file or directory")
}

func TestIncludesTestFilesInCommand(t *testing.T) {
	testConfig, ok := TestConfigurations["ruby"]
	assert.True(t, ok, "unexpectedly unable to find ruby test config")

	// this creates a config file in the test directory and removes it
	dir := filepath.Join(".", ".exercism")
	err := os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	f, err := os.Create(filepath.Join(dir, "config.json"))
	assert.NoError(t, err)

	_, err = f.WriteString(`{ "blurb": "Learn about the basics of Ruby by following a lasagna recipe.", "authors": ["iHiD", "pvcarrera"], "files": { "solution": ["lasagna.rb"], "test": ["lasagna_test.rb", "some_other_file.rb"], "exemplar": [".meta/exemplar.rb"] } } `)
	assert.NoError(t, err)

	cmd, err := testConfig.GetTestCommand()
	assert.NoError(t, err)
	assert.Equal(t, cmd, "ruby lasagna_test.rb some_other_file.rb")
}

func TestRustHasTrailingDashes(t *testing.T) {
	testConfig, ok := TestConfigurations["rust"]
	assert.True(t, ok, "unexpectedly unable to find rust test config")

	cmd, err := testConfig.GetTestCommand()
	assert.NoError(t, err)

	assert.True(t, strings.HasSuffix(cmd, "--"), "rust's test command should have trailing dashes")
}
