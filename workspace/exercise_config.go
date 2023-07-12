package workspace

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

const configFilename = "config.json"

var configFilepath = filepath.Join(ignoreSubdir, configFilename)

// ExerciseConfig contains exercise metadata.
// Note: we only use a subset of its fields
type ExerciseConfig struct {
	Files struct {
		Test []string `json:"test"`
	} `json:"files"`
}

// NewExerciseConfig reads exercise metadata from a file in the given directory.
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
