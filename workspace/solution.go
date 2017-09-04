package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/exercism/cli/visibility"
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
	Dir         string     `json:"-"`
}

// NewSolution reads solution metadata from a file in the given directory.
func NewSolution(dir string) (*Solution, error) {
	path := filepath.Join(dir, solutionFilename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return &Solution{}, err
	}
	var s Solution
	if err := json.Unmarshal(b, &s); err != nil {
		return &Solution{}, err
	}
	s.Dir = dir
	return &s, nil
}

// Suffix is the serial numeric value appended to an exercise directory.
// This is appended to avoid name conflicts, and does not indicate a particular
// iteration.
func (s *Solution) Suffix() string {
	return strings.Trim(strings.Replace(filepath.Base(s.Dir), s.Exercise, "", 1), "-.")
}

func (s *Solution) String() string {
	str := fmt.Sprintf("%s/%s", s.Track, s.Exercise)
	if s.Suffix() != "" {
		str = fmt.Sprintf("%s (%s)", str, s.Suffix())
	}
	if !s.IsRequester && s.Handle != "" {
		str = fmt.Sprintf("%s by @%s", str, s.Handle)
	}
	return str
}

// Write stores solution metadata to a file.
func (s *Solution) Write(dir string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, solutionFilename)
	if err := ioutil.WriteFile(path, b, os.FileMode(0644)); err != nil {
		return err
	}
	s.Dir = dir
	return visibility.HideFile(path)
}

// PathToParent is the relative path from the workspace to the parent dir.
func (s *Solution) PathToParent() string {
	var dir string
	if !s.IsRequester {
		dir = filepath.Join("users")
	}
	return filepath.Join(dir, s.Track)
}
