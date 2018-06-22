package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rgxSerialSuffix = regexp.MustCompile(`-\d*$`)

// Workspace represents a user's Exercism workspace.
// It may contain a user's own exercises, and other people's
// exercises that they've downloaded to look at or run locally.
type Workspace struct {
	Dir string
}

// New returns a configured workspace.
func New(dir string) Workspace {
	return Workspace{Dir: dir}
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

// SolutionPath returns the full path where the exercise will be stored.
// By default this the directory name matches that of the exercise, but if
// a different solution already exists, then a numeric suffix will be added
// to the name.
func (ws Workspace) SolutionPath(exercise, solutionID string) (string, error) {
	paths, err := ws.Locate(exercise)
	if !IsNotExist(err) && err != nil {
		return "", err
	}

	return ws.ResolveSolutionPath(paths, exercise, solutionID, IsSolutionPath)
}

// IsSolutionPath checks whether the given path contains the solution with the given ID.
func IsSolutionPath(solutionID, path string) (bool, error) {
	s, err := NewSolution(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return s.ID == solutionID, nil
}

// ResolveSolutionPath determines the path for the given exercise solution.
// It will locate an existing path, or indicate the name of a new path, if this is a new solution.
func (ws Workspace) ResolveSolutionPath(paths []string, exercise, solutionID string, existsFn func(string, string) (bool, error)) (string, error) {
	// Do we already have a directory for this solution?
	for _, path := range paths {
		ok, err := existsFn(solutionID, path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}
	}
	// If we didn't find the solution in one of the paths that
	// were passed in, we're going to construct some new ones
	// using a numeric suffix. Create a lookup table so we can
	// reject constructed paths if they match existing ones.
	m := map[string]bool{}
	for _, path := range paths {
		m[path] = true
	}
	suffix := 1
	root := filepath.Join(ws.Dir, exercise)
	path := root
	for {
		exists := m[path]
		if !exists {
			return path, nil
		}
		suffix++
		path = fmt.Sprintf("%s-%d", root, suffix)
	}
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
			return "", errors.New("couldn't find it")
		}
		if _, err := os.Lstat(path); os.IsNotExist(err) {
			return "", err
		}
		if _, err := os.Lstat(filepath.Join(path, solutionFilename)); err == nil {
			return path, nil
		}
		path = filepath.Dir(path)
	}
	return "", nil
}
