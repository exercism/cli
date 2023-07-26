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
	assert.Error(t, err)
	// any assertions about this error message have to work across all platforms, so be vague
	// unix: ".exercism/config.json: no such file or directory"
	// windows: "open .exercism\config.json: The system cannot find the path specified."
	assert.Contains(t, err.Error(), filepath.Join(".exercism", "config.json:"))
}

func TestIncludesSolutionAndTestFilesInCommand(t *testing.T) {
	testConfig, ok := TestConfigurations["prolog"]
	assert.True(t, ok, "unexpectedly unable to find prolog test config")

	// this creates a config file in the test directory and removes it
	dir := filepath.Join(".", ".exercism")
	defer os.RemoveAll(dir)
	err := os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, "config.json"))
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(`{ "blurb": "Learn about the basics of Prolog by following a lasagna recipe.", "authors": ["iHiD", "pvcarrera"], "files": { "solution": ["lasagna.pl"], "test": ["lasagna_tests.plt"] } } `)
	assert.NoError(t, err)

	cmd, err := testConfig.GetTestCommand()
	assert.NoError(t, err)
	assert.Equal(t, cmd, "swipl -f lasagna.pl -s lasagna_tests.plt -g run_tests,halt -t 'halt(1)'")
}

func TestIncludesTestFilesInCommand(t *testing.T) {
	testConfig, ok := TestConfigurations["ruby"]
	assert.True(t, ok, "unexpectedly unable to find ruby test config")

	// this creates a config file in the test directory and removes it
	dir := filepath.Join(".", ".exercism")
	defer os.RemoveAll(dir)
	err := os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, "config.json"))
	assert.NoError(t, err)
	defer f.Close()

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
