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

// Metadata contains metadata about a user's solution.
type Metadata struct {
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

// NewMetadata reads solution metadata from a file in the given directory.
func NewMetadata(dir string) (*Metadata, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, metadataFilepath))
	if err != nil {
		return &Metadata{}, err
	}
	var s Metadata
	if err := json.Unmarshal(b, &s); err != nil {
		return &Metadata{}, err
	}
	s.Dir = dir
	return &s, nil
}

// Suffix is the serial numeric value appended to an exercise directory.
// This is appended to avoid name conflicts, and does not indicate a particular
// iteration.
func (s *Metadata) Suffix() string {
	return strings.Trim(strings.Replace(filepath.Base(s.Dir), s.Exercise, "", 1), "-.")
}

func (s *Metadata) String() string {
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
func (s *Metadata) Write(dir string) error {
	b, err := json.Marshal(s)
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
	s.Dir = dir
	return nil
}

// PathToParent is the relative path from the workspace to the parent dir.
func (s *Metadata) PathToParent() string {
	var dir string
	if !s.IsRequester {
		dir = filepath.Join("users")
	}
	return filepath.Join(dir, s.Track)
}
