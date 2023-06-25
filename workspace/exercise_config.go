package workspace

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const configFilename = "config.json"

var configFilepath = filepath.Join(ignoreSubdir, configFilename)

// ExerciseConfig contains metadata about the problem.
// It's got a bunch more fields, but we don't need to read them
type ExerciseConfig struct {
	Files struct {
		Test []string `json:"test"`
	} `json:"files"`
}

// NewExerciseMetadata reads exercise metadata from a file in the given directory.
func NewExerciseConfig(dir string) (*ExerciseConfig, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, configFilepath))
	if err != nil {
		return nil, err
	}
	var config ExerciseConfig
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
