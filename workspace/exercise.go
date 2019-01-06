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
	return filepath.Join(e.Filepath(), legacyMetadataFilename)
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

// HasLegacyMetadata checks for the presence of a legacy exercise metadata file.
// If there is no such file, it could also be an unrelated directory.
func (e Exercise) HasLegacyMetadata() (bool, error) {
	_, err := os.Lstat(e.LegacyMetadataFilepath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

// MigrationStatus represents the result of migrating a legacy metadata file.
type MigrationStatus int

// MigrationStatus
const (
	MigrationStatusNoop MigrationStatus = iota
	MigrationStatusMigrated
	MigrationStatusRemoved
)

func (m MigrationStatus) String() string {
	switch m {
	case MigrationStatusMigrated:
		return "\nMigrated metadata\n"
	case MigrationStatusRemoved:
		return "\nRemoved legacy metadata\n"
	default:
		return ""
	}
}

// MigrateLegacyMetadataFile migrates a legacy metadata file to the modern location.
// This is a noop if the metadata file isn't legacy.
// If both legacy and modern metadata files exist, the legacy file will be deleted.
func (e Exercise) MigrateLegacyMetadataFile() (MigrationStatus, error) {
	if ok, _ := e.HasLegacyMetadata(); !ok {
		return MigrationStatusNoop, nil
	}
	if err := os.MkdirAll(filepath.Dir(e.MetadataFilepath()), os.FileMode(0755)); err != nil {
		return MigrationStatusNoop, err
	}
	if ok, _ := e.HasMetadata(); !ok {
		if err := os.Rename(e.LegacyMetadataFilepath(), e.MetadataFilepath()); err != nil {
			return MigrationStatusNoop, err
		}
		return MigrationStatusMigrated, nil
	}
	if err := os.Remove(e.LegacyMetadataFilepath()); err != nil {
		return MigrationStatusNoop, err
	}
	return MigrationStatusRemoved, nil
}
