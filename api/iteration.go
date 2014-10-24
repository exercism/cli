package api

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	msgUnidentifiable = "unable to identify track and problem"
)

// Iteration represents a version of a particular exercise.
// This gets submitted to the API.
type Iteration struct {
	Key      string `json:"key"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Dir      string `json:"dir"`
	File     string `json:"-"`
	Language string `json:"-"`
	Problem  string `json:"-"`
}

// RelativePath returns the path relative to the exercism dir.
func (iter *Iteration) RelativePath() string {
	if iter.Path != "" {
		return iter.Path
	}

	if len(iter.Dir) > len(iter.File) {
		return ""
	}
	iter.Path = iter.File[len(iter.Dir):]
	return iter.Path
}

// Identify attempts to determine the track and problem of an iteration.
func (iter *Iteration) Identify() error {
	if !strings.HasPrefix(strings.ToLower(iter.File), strings.ToLower(iter.Dir)) {
		return fmt.Errorf(msgUnidentifiable)
	}

	segments := strings.Split(iter.RelativePath(), string(filepath.Separator))
	// file is always the absolute path, so the first segment will be empty
	if len(segments) < 4 {
		return fmt.Errorf(msgUnidentifiable)
	}

	iter.Language = segments[1]
	iter.Problem = segments[2]
	return nil
}
