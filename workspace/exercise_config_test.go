package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExerciseConfig(t *testing.T) {
	dir, err := os.MkdirTemp("", "exercise_config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.Mkdir(filepath.Join(dir, ".exercism"), os.ModePerm)
	assert.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, ".exercism", "config.json"))
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(`{ "blurb": "Learn about the basics of Ruby by following a lasagna recipe.", "authors": ["iHiD", "pvcarrera"], "files": { "solution": ["lasagna.rb"], "test": ["lasagna_test.rb"], "exemplar": [".meta/exemplar.rb"] } } `)
	assert.NoError(t, err)

	ec, err := NewExerciseConfig(dir)
	assert.NoError(t, err)

	assert.Equal(t, ec.Files.Solution, []string{"lasagna.rb"})
	solutionFiles, err := ec.GetSolutionFiles()
	assert.NoError(t, err)
	assert.Equal(t, solutionFiles, []string{"lasagna.rb"})

	assert.Equal(t, ec.Files.Test, []string{"lasagna_test.rb"})
	testFiles, err := ec.GetTestFiles()
	assert.NoError(t, err)
	assert.Equal(t, testFiles, []string{"lasagna_test.rb"})
}

func TestExerciseConfigNoTestKey(t *testing.T) {
	dir, err := os.MkdirTemp("", "exercise_config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.Mkdir(filepath.Join(dir, ".exercism"), os.ModePerm)
	assert.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, ".exercism", "config.json"))
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(`{ "blurb": "Learn about the basics of Ruby by following a lasagna recipe.", "authors": ["iHiD", "pvcarrera"], "files": { "exemplar": [".meta/exemplar.rb"] } } `)
	assert.NoError(t, err)

	ec, err := NewExerciseConfig(dir)
	assert.NoError(t, err)

	_, err = ec.GetSolutionFiles()
	assert.Error(t, err, "no `files.solution` key in your `config.json`")
	_, err = ec.GetTestFiles()
	assert.Error(t, err, "no `files.test` key in your `config.json`")
}

func TestMissingExerciseConfig(t *testing.T) {
	dir, err := os.MkdirTemp("", "exercise_config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	_, err = NewExerciseConfig(dir)
	assert.Error(t, err)
	// any assertions about this error message have to work across all platforms, so be vague
	// unix: ".exercism/config.json: no such file or directory"
	// windows: "open .exercism\config.json: The system cannot find the path specified."
	assert.Contains(t, err.Error(), filepath.Join(".exercism", "config.json:"))
}

func TestInvalidExerciseConfig(t *testing.T) {
	dir, err := os.MkdirTemp("", "exercise_config")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.Mkdir(filepath.Join(dir, ".exercism"), os.ModePerm)
	assert.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, ".exercism", "config.json"))
	assert.NoError(t, err)
	defer f.Close()

	// invalid JSON
	_, err = f.WriteString(`{ "blurb": "Learn about the basics of Ruby by following a lasagna recipe.", "authors": ["iHiD", "pvcarr `)
	assert.NoError(t, err)

	_, err = NewExerciseConfig(dir)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unexpected end of JSON input"))
}
