package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const configFilename = "config.json"

var configFilepath = filepath.Join(ignoreSubdir, configFilename)

// ExerciseConfig contains exercise metadata.
// Note: we only use a subset of its fields
type ExerciseConfig struct {
	Files struct {
		Solution []string `json:"solution"`
		Test     []string `json:"test"`
	} `json:"files"`
}

// NewExerciseConfig reads exercise metadata from a file in the given directory.
func NewExerciseConfig(dir string) (*ExerciseConfig, error) {
	b, err := os.ReadFile(filepath.Join(dir, configFilepath))
	if err != nil {
		return nil, err
	}
	var config ExerciseConfig
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetTestFiles finds returns the names of the file(s) that hold unit tests for this exercise, if any
func (c *ExerciseConfig) GetSolutionFiles() ([]string, error) {
	result := c.Files.Solution
	if result == nil {
		// solution file(s) key was missing in config json, which is an error when calling this fuction
		return []string{}, errors.New("no `files.solution` key in your `config.json`. Was it removed by mistake?")
	}

	return result, nil
}

// GetTestFiles finds returns the names of the file(s) that hold unit tests for this exercise, if any
func (c *ExerciseConfig) GetTestFiles() ([]string, error) {
	result := c.Files.Test
	if result == nil {
		// test file(s) key was missing in config json, which is an error when calling this fuction
		return []string{}, errors.New("no `files.test` key in your `config.json`. Was it removed by mistake?")
	}

	return result, nil
}
