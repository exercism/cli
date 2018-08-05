package workspace

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var errMissingMetadata = errors.New("no solution metadata file found")

// IsMissingMetadata verifies the type of error.
func IsMissingMetadata(err error) bool {
	return err == errMissingMetadata
}

var rgxSerialSuffix = regexp.MustCompile(`-\d*$`)

// Workspace represents a user's Exercism workspace.
// It may contain a user's own exercises, and other people's
// exercises that they've downloaded to look at or run locally.
type Workspace struct {
	Dir string
}

// New returns a configured workspace.
func New(dir string) (Workspace, error) {
	_, err := os.Lstat(dir)
	if err != nil {
		return Workspace{}, err
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{Dir: dir}, nil
}

// PotentialExercises are a first-level guess at the user's exercises.
// It looks at the workspace structurally, and guesses based on
// the location of the directory. E.g. any top level directory
// within the workspace (except 'users') is assumed to be a
// track, and any directory within there again is assumed to
// be an exercise.
func (ws Workspace) PotentialExercises() ([]Exercise, error) {
	exercises := []Exercise{}

	topInfos, err := ioutil.ReadDir(ws.Dir)
	if err != nil {
		return nil, err
	}
	for _, topInfo := range topInfos {
		if !topInfo.IsDir() {
			continue
		}

		if topInfo.Name() == "users" {
			continue
		}

		subInfos, err := ioutil.ReadDir(filepath.Join(ws.Dir, topInfo.Name()))
		if err != nil {
			return nil, err
		}

		for _, subInfo := range subInfos {
			if !subInfo.IsDir() {
				continue
			}

			exercises = append(exercises, Exercise{Track: topInfo.Name(), Slug: subInfo.Name(), Root: ws.Dir})
		}
	}

	return exercises, nil
}

// Exercises returns the user's exercises within the workspace.
// This doesn't find legacy exercises where the metadata is missing.
func (ws Workspace) Exercises() ([]Exercise, error) {
	candidates, err := ws.PotentialExercises()
	if err != nil {
		return nil, err
	}

	exercises := make([]Exercise, 0, len(candidates))
	for _, candidate := range candidates {
		ok, err := candidate.HasMetadata()
		if err != nil {
			return nil, err
		}
		if ok {
			exercises = append(exercises, candidate)
		}
	}
	return exercises, nil
}

// Locate the matching directories within the workspace.
// This will look for an exact match on absolute or relative paths.
// If given the base name of a directory with no path information it
// It will look for all directories with that name, or that are
// named with a numerical suffix.
func (ws Workspace) Locate(exercise string) ([]string, error) {
	// First assume it's a path.
	dir := exercise

	// If it's not an absolute path, make it one.
	if !filepath.IsAbs(dir) {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
	}

	// If it exists, we were right. It's a path.
	if _, err := os.Stat(dir); err == nil {
		if !strings.HasPrefix(dir, ws.Dir) {
			return nil, ErrNotInWorkspace(exercise)
		}

		src, err := filepath.EvalSymlinks(dir)
		if err == nil {
			return []string{src}, nil
		}
	}

	// If the argument is a path, then we should have found it by now.
	if strings.Contains(exercise, string(os.PathSeparator)) {
		return nil, ErrNotExist(exercise)
	}

	var paths []string
	// Look through the entire workspace tree to find any matches.
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a symlink, follow it, then get the file info of the target.
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			src, err := filepath.EvalSymlinks(path)
			if err == nil {
				path = src
			}
			info, err = os.Lstat(path)
			if err != nil {
				return err
			}
		}

		if !info.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), exercise) {
			// We're trying to find any directories that match either the exact name
			// or the name with a numeric suffix.
			// E.g. if passed 'bat', then we should match 'bat', 'bat-2', 'bat-200',
			// but not 'batten'.
			suffix := strings.Replace(filepath.Base(path), exercise, "", 1)
			if suffix == "" || rgxSerialSuffix.MatchString(suffix) {
				paths = append(paths, path)
			}
		}
		return nil
	}

	// If the workspace directory is a symlink, resolve that first.
	root := ws.Dir
	src, err := filepath.EvalSymlinks(root)
	if err == nil {
		root = src
	}

	filepath.Walk(root, walkFn)

	if len(paths) == 0 {
		return nil, ErrNotExist(exercise)
	}
	return paths, nil
}

// SolutionDir determines the root directory of a solution.
// This is the directory that contains the solution metadata file.
func (ws Workspace) SolutionDir(s string) (string, error) {
	if !strings.HasPrefix(s, ws.Dir) {
		return "", errors.New("not in workspace")
	}

	path := s
	for {
		if path == ws.Dir {
			return "", errMissingMetadata
		}
		if _, err := os.Lstat(path); os.IsNotExist(err) {
			return "", err
		}
		if _, err := os.Lstat(filepath.Join(path, solutionFilename)); err == nil {
			return path, nil
		}
		path = filepath.Dir(path)
	}
}
