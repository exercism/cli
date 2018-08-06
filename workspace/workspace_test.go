package workspace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspacePotentialExercises(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "walk")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	a1 := filepath.Join(tmpDir, "track-a", "exercise-one")
	b1 := filepath.Join(tmpDir, "track-b", "exercise-one")
	b2 := filepath.Join(tmpDir, "track-b", "exercise-two")

	// It should ignore other people's exercises.
	alice := filepath.Join(tmpDir, "users", "alice", "track-a", "exercise-one")

	// It should ignore nested dirs within exercises.
	nested := filepath.Join(a1, "subdir", "deeper-dir", "another-deep-dir")

	for _, path := range []string{a1, b1, b2, alice, nested} {
		err := os.MkdirAll(path, os.FileMode(0755))
		assert.NoError(t, err)
	}

	ws, err := New(tmpDir)
	assert.NoError(t, err)

	exercises, err := ws.PotentialExercises()
	assert.NoError(t, err)
	if assert.Equal(t, 3, len(exercises)) {
		paths := make([]string, len(exercises))
		for i, e := range exercises {
			paths[i] = e.Path()
		}

		sort.Strings(paths)
		assert.Equal(t, paths[0], "track-a/exercise-one")
		assert.Equal(t, paths[1], "track-b/exercise-one")
		assert.Equal(t, paths[2], "track-b/exercise-two")
	}
}

func TestWorkspaceExercises(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "walk-with-metadata")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	a1 := filepath.Join(tmpDir, "track-a", "exercise-one")
	a2 := filepath.Join(tmpDir, "track-a", "exercise-two") // no metadata
	b1 := filepath.Join(tmpDir, "track-b", "exercise-one")
	b2 := filepath.Join(tmpDir, "track-b", "exercise-two")

	for _, path := range []string{a1, a2, b1, b2} {
		err := os.MkdirAll(path, os.FileMode(0755))
		assert.NoError(t, err)

		if path != a2 {
			err = ioutil.WriteFile(filepath.Join(path, metadataFilename), []byte{}, os.FileMode(0600))
			assert.NoError(t, err)
		}
	}

	ws, err := New(tmpDir)
	assert.NoError(t, err)

	exercises, err := ws.Exercises()
	assert.NoError(t, err)
	if assert.Equal(t, 3, len(exercises)) {
		paths := make([]string, len(exercises))
		for i, e := range exercises {
			paths[i] = e.Path()
		}

		sort.Strings(paths)
		assert.Equal(t, paths[0], "track-a/exercise-one")
		assert.Equal(t, paths[1], "track-b/exercise-one")
		assert.Equal(t, paths[2], "track-b/exercise-two")
	}
}

func TestSolutionPath(t *testing.T) {
	root := filepath.Join("..", "fixtures", "solution-path", "creatures")
	ws, err := New(root)
	assert.NoError(t, err)

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
	tmpDir, err := ioutil.TempDir("", "resolve-solution-path")
	defer os.RemoveAll(tmpDir)
	ws, err := New(tmpDir)
	assert.NoError(t, err)

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

	ws, err := New(filepath.Join(root, "workspace"))
	assert.NoError(t, err)

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
