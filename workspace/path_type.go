package workspace

import (
	"os"
	"path/filepath"
)

// PathType is either a path to a dir or file, or the name of an exercise.
type PathType int

const (
	// TypeExerciseID is the name of an exercise.
	TypeExerciseID PathType = iota
	// TypeDir is a relative or absolute path to a directory.
	TypeDir
	// TypeFile is a relative or absolute path to a file.
	TypeFile
)

// DetectPathType determines whether the given path is a directory, a file, or the name of an exercise.
func DetectPathType(path string) (PathType, error) {
	// If it's not an absolute path, make it one.
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return -1, err
		}
	}

	// If it doesn't exist, then it's an exercise name.
	// We'll have to walk the workspace to find it.
	if _, err := os.Stat(path); err != nil {
		return TypeExerciseID, nil
	}

	// We found it. It's an actual path of some sort.
	info, err := os.Lstat(path)
	if err != nil {
		return -1, err
	}

	// If it's a symlink, resolve it.
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		src, err := filepath.EvalSymlinks(path)
		if err != nil {
			return -1, err
		}
		path = src
		// Overwrite the symlinked info with the source info.
		info, err = os.Lstat(path)
		if err != nil {
			return -1, err
		}
	}

	if info.IsDir() {
		return TypeDir, nil
	}
	return TypeFile, nil
}
