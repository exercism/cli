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

const metadataFilename = ".solution.json"

// Metadata contains metadata about a user's exercise.
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

// NewMetadata reads exercise metadata from a file in the given directory.
func NewMetadata(dir string) (*Metadata, error) {
	path := filepath.Join(dir, metadataFilename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var metadata Metadata
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, err
	}
	metadata.Dir = dir
	return &metadata, nil
}

// Suffix is the serial numeric value appended to an exercise directory.
// This is appended to avoid name conflicts, and does not indicate a particular
// iteration.
func (metadata *Metadata) Suffix() string {
	return strings.Trim(strings.Replace(filepath.Base(metadata.Dir), metadata.Exercise, "", 1), "-.")
}

func (metadata *Metadata) String() string {
	str := fmt.Sprintf("%s/%s", metadata.Track, metadata.Exercise)
	if metadata.Suffix() != "" {
		str = fmt.Sprintf("%s (%s)", str, metadata.Suffix())
	}
	if !metadata.IsRequester && metadata.Handle != "" {
		str = fmt.Sprintf("%s by @%s", str, metadata.Handle)
	}
	return str
}

// Write stores solution metadata to a file.
func (metadata *Metadata) Write(dir string) error {
	b, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, metadataFilename)

	// Hack because ioutil.WriteFile fails on hidden files
	visibility.ShowFile(path)

	if err := ioutil.WriteFile(path, b, os.FileMode(0600)); err != nil {
		return err
	}
	metadata.Dir = dir
	return visibility.HideFile(path)
}

// PathToParent is the relative path from the workspace to the parent dir.
func (metadata *Metadata) PathToParent() string {
	var dir string
	if !metadata.IsRequester {
		dir = filepath.Join("users")
	}
	return filepath.Join(dir, metadata.Track)
}
