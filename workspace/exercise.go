package workspace

import (
	"path"
	"path/filepath"
)

// Exercise is an implementation of a problem in a track.
type Exercise struct {
	Root  string
	Track string
	Slug  string
}

// Path is the normalized relative path.
// It always has forward slashes, regardless
// of the operating system.
func (e Exercise) Path() string {
	return path.Join(e.Track, e.Slug)
}

// Filepath is the absolute path on the filesystem.
func (e Exercise) Filepath() string {
	return filepath.Join(e.Root, e.Track, e.Slug)
}
