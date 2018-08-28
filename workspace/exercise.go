package workspace

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// Exercise is an implementation of a problem in a track.
type Exercise struct {
	Root      string
	Track     string
	Slug      string
	Documents []Document
}

// NewExerciseFromDir constructs an exercise given the exercise directory.
func NewExerciseFromDir(dir string) Exercise {
	slug := filepath.Base(dir)
	dir = filepath.Dir(dir)
	track := filepath.Base(dir)
	root := filepath.Dir(dir)
	return Exercise{Root: root, Track: track, Slug: slug}
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
	return filepath.Join(e.Filepath(), metadataFilepath)
}

// LegacyMetadataFilepath is the absolute path to the legacy exercise metadata.
func (e Exercise) LegacyMetadataFilepath() string {
	return filepath.Join(e.Filepath(), legacySolutionFilename)
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

// MigrateLegacyMetadataFile migrates a legacy metadata to the modern location.
// This is a noop if the metadata file isn't legacy.
// If both legacy and modern metadata files exist, the legacy file will be deleted.
func (e Exercise) MigrateLegacyMetadataFile() error {
	legacyMetadataFilepath := e.LegacyMetadataFilepath()
	metadataFilepath := e.MetadataFilepath()

	if _, err := os.Lstat(legacyMetadataFilepath); err != nil {
		return nil
	}
	if err := createIgnoreSubdir(filepath.Dir(legacyMetadataFilepath)); err != nil {
		return err
	}
	if _, err := os.Lstat(metadataFilepath); err != nil {
		if err := os.Rename(legacyMetadataFilepath, metadataFilepath); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\nMigrated metadata to %s\n", metadataFilepath)
	} else {
		if err := os.Remove(legacyMetadataFilepath); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\nRemoved legacy metadata: %s\n", legacyMetadataFilepath)
	}
	return nil
}
