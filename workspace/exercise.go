package workspace

import (
	"os"
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

// MetadataFilepath is the absolute path to the exercise metadata.
func (e Exercise) MetadataFilepath() string {
	return filepath.Join(e.Filepath(), solutionFilename)
}

// MetadataDir returns the directory that the exercise metadata lives in.
// For now this is the exercise directory.
func (e Exercise) MetadataDir() string {
	return e.Filepath()
}

// HasMetadata checks for the presence of an exercise metadata file.
// If there is no such file, this may be a legacy exercise.
// It could also be an unrelated directory.
func (e Exercise) HasMetadata() (bool, error) {
	_, err := os.Lstat(e.MetadataFilepath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}
