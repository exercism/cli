package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/exercism/cli/visibility"
)

const IgnoreSubdir = ".exercism"
const SolutionFilename = "solution.json"

var solutionRelPath = filepath.Join(IgnoreSubdir, SolutionFilename)

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
	AutoApprove bool       `json:"auto_approve"`
}

// NewSolution reads solution metadata from a file in the given directory.
func NewSolution(dir string) (*Solution, error) {
	path := filepath.Join(dir, solutionRelPath)
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

	path, err := createIgnoreSubdir(dir)
	if err != nil {
		return err
	}

	// Hack because ioutil.WriteFile fails on hidden files
	visibility.ShowFile(path)

	path = filepath.Join(dir, solutionRelPath)
	if err := ioutil.WriteFile(path, b, os.FileMode(0600)); err != nil {
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

func createIgnoreSubdir(path string) (string, error) {
	path = filepath.Join(path, IgnoreSubdir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.FileMode(0755)); err != nil {
			return "", fmt.Errorf("failed to create directory: %s", path)
		}
	}
	return path, nil
}

func migrateLegacySolutionFile(legacySolutionPath string, solutionPath string) error {
	if _, err := createIgnoreSubdir(filepath.Dir(legacySolutionPath)); err != nil {
		return err
	}
	if err := os.Rename(legacySolutionPath, solutionPath); err != nil {
		return errors.New("failed migrating legacy solution file")
	}
	return nil
}
