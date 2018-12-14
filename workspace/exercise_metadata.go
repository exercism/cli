package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const metadataFilename = "metadata.json"
const legacyMetadataFilename = ".solution.json"
const ignoreSubdir = ".exercism"

var metadataFilepath = filepath.Join(ignoreSubdir, metadataFilename)

// ExerciseMetadata contains metadata about a user's exercise.
type ExerciseMetadata struct {
	Track       string     `json:"track"`
	Exercise    string     `json:"exercise"`
	ID          string     `json:"id"`
	Team        string     `json:"team,omitempty"`
	URL         string     `json:"url"`
	Handle      string     `json:"handle"`
	IsRequester bool       `json:"is_requester"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Dir         string     `json:"-"`
	AutoApprove bool       `json:"auto_approve"`
}

// NewExerciseMetadata reads exercise metadata from a file in the given directory.
func NewExerciseMetadata(dir string) (*ExerciseMetadata, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, metadataFilepath))
	if err != nil {
		return nil, err
	}
	var metadata ExerciseMetadata
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, err
	}
	metadata.Dir = dir
	return &metadata, nil
}

// Suffix is the serial numeric value appended to an exercise directory.
// This is appended to avoid name conflicts, and does not indicate a particular
// iteration.
func (em *ExerciseMetadata) Suffix() string {
	return strings.Trim(strings.Replace(filepath.Base(em.Dir), em.Exercise, "", 1), "-.")
}

func (em *ExerciseMetadata) String() string {
	str := fmt.Sprintf("%s/%s", em.Track, em.Exercise)
	if em.Suffix() != "" {
		str = fmt.Sprintf("%s (%s)", str, em.Suffix())
	}
	if !em.IsRequester && em.Handle != "" {
		str = fmt.Sprintf("%s by @%s", str, em.Handle)
	}
	return str
}

// Write stores exercise metadata to a file.
func (em ExerciseMetadata) Write(dir string) error {
	b, err := json.Marshal(em)
	if err != nil {
		return err
	}
	metadataAbsoluteFilepath := filepath.Join(dir, metadataFilepath)
	if err = os.MkdirAll(filepath.Dir(metadataAbsoluteFilepath), os.FileMode(0755)); err != nil {
		return err
	}
	if err = ioutil.WriteFile(metadataAbsoluteFilepath, b, os.FileMode(0600)); err != nil {
		return err
	}
	em.Dir = dir
	return nil
}

// PathToParent is the relative path from the workspace to the parent dir.
func (em *ExerciseMetadata) PathToParent() string {
	var dir string
	if !em.IsRequester {
		dir = filepath.Join("users")
	}
	return filepath.Join(dir, em.Track)
}
