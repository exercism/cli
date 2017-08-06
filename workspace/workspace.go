package workspace

import (
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
