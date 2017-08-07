package workspace

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const solutionFilename = ".solution.json"

// Solution contains metadata about a user's solution.
type Solution struct {
	Track       string     `json:"track"`
	Exercise    string     `json:"exercise"`
	ID          string     `json:"id"`
	URL         string     `json:"url"`
	Handle      string     `json:"handle"`
	IsRequester bool       `json:"is_requester"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}

// NewSolution reads solution metadata from a file in the given directory.
func NewSolution(dir string) (Solution, error) {
	path := filepath.Join(dir, solutionFilename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Solution{}, err
	}
	var s Solution
	if err := json.Unmarshal(b, &s); err != nil {
		return Solution{}, err
	}
	return s, nil
}

// Write stores solution metadata to a file.
func (s Solution) Write(dir string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, solutionFilename)
	return ioutil.WriteFile(path, b, os.FileMode(0644))
}
