package workspace

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateErrors(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	ws := New(filepath.Join(root, "workspace"))

	tests := []struct {
		desc, arg string
		errFn     func(error) bool
	}{
		{
			desc:  "absolute path outside of workspace",
			arg:   filepath.Join(root, "equipment", "bat"),
			errFn: IsNotInWorkspace,
		},
		{
			desc:  "absolute path in workspace not found",
			arg:   filepath.Join(ws.Dir, "creatures", "pig"),
			errFn: IsNotExist,
		},
		{
			desc:  "relative path is outside of workspace",
			arg:   filepath.Join("..", "fixtures", "locate-exercise", "equipment", "bat"),
			errFn: IsNotInWorkspace,
		},
		{
			desc:  "relative path in workspace not found",
			arg:   filepath.Join("..", "fixtures", "locate-exercise", "workspace", "creatures", "pig"),
			errFn: IsNotExist,
		},
		{
			desc:  "exercise name not found in workspace",
			arg:   "pig",
			errFn: IsNotExist,
		},
	}

	for _, test := range tests {
		_, err := ws.Locate(test.arg)
		assert.True(t, test.errFn(err), fmt.Sprintf("test: %s (arg: %s), %#v", test.desc, test.arg, err))
	}
}

func TestLocate(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	wsPrimary := New(filepath.Join(root, "workspace"))
	wsSymbolic := New(filepath.Join(root, "symlinked-workspace"))

	tests := []struct {
		desc      string
		workspace Workspace
		in        string
		out       []string
	}{
		{
			desc:      "find absolute path within workspace",
			workspace: wsPrimary,
			in:        filepath.Join(wsPrimary.Dir, "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find absolute path within symlinked workspace",
			workspace: wsSymbolic,
			in:        filepath.Join(wsSymbolic.Dir, "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find relative path within workspace",
			workspace: wsPrimary,
			in:        filepath.Join("..", "fixtures", "locate-exercise", "workspace", "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in default location",
			workspace: wsPrimary,
			in:        "horse",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in a symlinked workspace",
			workspace: wsSymbolic,
			in:        "horse",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in a subtree",
			workspace: wsPrimary,
			in:        "fly",
			out:       []string{filepath.Join(wsPrimary.Dir, "friends", "alice", "creatures", "fly")},
		},
		{
			desc:      "don't be confused by a file named the same as an exercise",
			workspace: wsPrimary,
			in:        "duck",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "duck")},
		},
		{
			desc:      "don't be confused by a symlinked file named the same as an exercise",
			workspace: wsPrimary,
			in:        "date",
			out:       []string{filepath.Join(wsPrimary.Dir, "actions", "date")},
		},
		{
			desc:      "find all the exercises with the same name",
			workspace: wsPrimary,
			in:        "bat",
			out: []string{
				filepath.Join(wsPrimary.Dir, "creatures", "bat"),
				filepath.Join(wsPrimary.Dir, "friends", "alice", "creatures", "bat"),
			},
		},
		{
			desc:      "find copies of exercise with suffix",
			workspace: wsPrimary,
			in:        "crane",
			out: []string{
				filepath.Join(wsPrimary.Dir, "creatures", "crane"),
				filepath.Join(wsPrimary.Dir, "creatures", "crane-2"),
			},
		},
		{
			desc:      "find exercises that are symlinks",
			workspace: wsPrimary,
			in:        "squash",
			out: []string{
				filepath.Join(wsPrimary.Dir, "..", "food", "squash"),
				filepath.Join(wsPrimary.Dir, "actions", "squash"),
			},
		},
	}

	for _, test := range tests {
		dirs, err := test.workspace.Locate(test.in)

		sort.Strings(dirs)
		sort.Strings(test.out)

		assert.NoError(t, err, test.desc)
		assert.Equal(t, test.out, dirs, test.desc)
	}
}

func TestSolutionPath(t *testing.T) {
	root := filepath.Join("..", "fixtures", "solution-path", "creatures")
	ws := New(root)

	// An existing exercise.
	path, err := ws.SolutionPath("gazelle", "ccc")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "gazelle-3"), path)

	path, err = ws.SolutionPath("gazelle", "abc")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "gazelle-4"), path)

	// A new exercise.
	path, err = ws.SolutionPath("lizard", "abc")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "lizard"), path)
}

func TestIsSolutionPath(t *testing.T) {
	root := filepath.Join("..", "fixtures", "is-solution-path")

	ok, err := IsSolutionPath("abc", filepath.Join(root, "yepp"))
	assert.NoError(t, err)
	assert.True(t, ok)

	// The ID has to actually match.
	ok, err = IsSolutionPath("xxx", filepath.Join(root, "yepp"))
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = IsSolutionPath("abc", filepath.Join(root, "nope"))
	assert.NoError(t, err)
	assert.False(t, ok)

	_, err = IsSolutionPath("abc", filepath.Join(root, "broken"))
	assert.Error(t, err)
}

func TestResolveSolutionPath(t *testing.T) {
	ws := New("tmp")

	existsFn := func(solutionID, path string) (bool, error) {
		pathToSolutionID := map[string]string{
			filepath.Join(ws.Dir, "pig"):      "xxx",
			filepath.Join(ws.Dir, "gecko"):    "aaa",
			filepath.Join(ws.Dir, "gecko-2"):  "xxx",
			filepath.Join(ws.Dir, "gecko-3"):  "ccc",
			filepath.Join(ws.Dir, "bat"):      "aaa",
			filepath.Join(ws.Dir, "dog"):      "aaa",
			filepath.Join(ws.Dir, "dog-2"):    "bbb",
			filepath.Join(ws.Dir, "dog-3"):    "ccc",
			filepath.Join(ws.Dir, "rabbit"):   "aaa",
			filepath.Join(ws.Dir, "rabbit-2"): "bbb",
			filepath.Join(ws.Dir, "rabbit-4"): "ccc",
		}
		return pathToSolutionID[path] == solutionID, nil
	}

	tests := []struct {
		desc     string
		paths    []string
		exercise string
		expected string
	}{
		{
			desc:     "If we don't have that exercise yet, it gets the default name.",
			exercise: "duck",
			paths:    []string{},
			expected: filepath.Join(ws.Dir, "duck"),
		},
		{
			desc:     "If we already have a directory for the solution in question, return it.",
			exercise: "pig",
			paths: []string{
				filepath.Join(ws.Dir, "pig"),
			},
			expected: filepath.Join(ws.Dir, "pig"),
		},
		{
			desc:     "If we already have multiple solutions, and this is one of them, find it.",
			exercise: "gecko",
			paths: []string{
				filepath.Join(ws.Dir, "gecko"),
				filepath.Join(ws.Dir, "gecko-2"),
				filepath.Join(ws.Dir, "gecko-3"),
			},
			expected: filepath.Join(ws.Dir, "gecko-2"),
		},
		{
			desc:     "If we already have a solution, but this is a new one, add a suffix.",
			exercise: "bat",
			paths: []string{
				filepath.Join(ws.Dir, "bat"),
			},
			expected: filepath.Join(ws.Dir, "bat-2"),
		},
		{
			desc:     "If we already have multiple solutions, but this is a new one, add a new suffix.",
			exercise: "dog",
			paths: []string{
				filepath.Join(ws.Dir, "dog"),
				filepath.Join(ws.Dir, "dog-2"),
				filepath.Join(ws.Dir, "dog-3"),
			},
			expected: filepath.Join(ws.Dir, "dog-4"),
		},
		{
			desc:     "Use the first available suffix.",
			exercise: "rabbit",
			paths: []string{
				filepath.Join(ws.Dir, "rabbit"),
				filepath.Join(ws.Dir, "rabbit-2"),
				filepath.Join(ws.Dir, "rabbit-4"),
			},
			expected: filepath.Join(ws.Dir, "rabbit-3"),
		},
	}

	for _, test := range tests {
		path, err := ws.ResolveSolutionPath(test.paths, test.exercise, "xxx", existsFn)
		assert.NoError(t, err, test.desc)
		assert.Equal(t, test.expected, path, test.desc)
	}
}

func TestSolutionDir(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solution-dir")

	ws := New(filepath.Join(root, "workspace"))

	tests := []struct {
		path string
		ok   bool
	}{
		{
			path: filepath.Join(ws.Dir, "exercise"),
			ok:   true,
		},
		{
			path: filepath.Join(ws.Dir, "exercise", "file.txt"),
			ok:   true,
		},
		{
			path: filepath.Join(ws.Dir, "exercise", "in", "a", "subdir", "file.txt"),
			ok:   true,
		},
		{
			path: filepath.Join(ws.Dir, "exercise", "in", "a"),
			ok:   true,
		},
		{
			path: filepath.Join(ws.Dir, "not-exercise", "file.txt"),
			ok:   false,
		},
		{
			path: filepath.Join(root, "file.txt"),
			ok:   false,
		},
		{
			path: filepath.Join(ws.Dir, "exercise", "no-such-file.txt"),
			ok:   false,
		},
	}

	for _, test := range tests {
		dir, err := ws.SolutionDir(test.path)
		if !test.ok {
			assert.Error(t, err, test.path)
			continue
		}

		assert.Equal(t, filepath.Join(ws.Dir, "exercise"), dir, test.path)
	}
}
