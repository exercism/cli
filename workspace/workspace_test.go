package workspace

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspacePotentialExercises(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "walk")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	a1 := filepath.Join(tmpDir, "track-a", "exercise-one")
	b1 := filepath.Join(tmpDir, "track-b", "exercise-one")
	b2 := filepath.Join(tmpDir, "track-b", "exercise-two")

	// It should find teams exercises
	team := filepath.Join(tmpDir, "teams", "some-team", "track-c", "exercise-one")

	// It should ignore other people's exercises.
	alice := filepath.Join(tmpDir, "users", "alice", "track-a", "exercise-one")

	// It should ignore nested dirs within exercises.
	nested := filepath.Join(a1, "subdir", "deeper-dir", "another-deep-dir")

	for _, path := range []string{a1, b1, b2, team, alice, nested} {
		err := os.MkdirAll(path, os.FileMode(0755))
		assert.NoError(t, err)
	}

	ws, err := New(tmpDir)
	assert.NoError(t, err)

	exercises, err := ws.PotentialExercises()
	assert.NoError(t, err)
	if assert.Equal(t, 4, len(exercises)) {
		paths := make([]string, len(exercises))
		for i, e := range exercises {
			paths[i] = e.Path()
		}

		sort.Strings(paths)
		assert.Equal(t, paths[0], "track-a/exercise-one")
		assert.Equal(t, paths[1], "track-b/exercise-one")
		assert.Equal(t, paths[2], "track-b/exercise-two")
		assert.Equal(t, paths[3], "track-c/exercise-one")
	}
}

func TestWorkspaceExercises(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "walk-with-metadata")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	a1 := filepath.Join(tmpDir, "track-a", "exercise-one")
	a2 := filepath.Join(tmpDir, "track-a", "exercise-two") // no metadata
	b1 := filepath.Join(tmpDir, "track-b", "exercise-one")
	b2 := filepath.Join(tmpDir, "track-b", "exercise-two")

	for _, path := range []string{a1, a2, b1, b2} {
		metadataAbsoluteFilepath := filepath.Join(path, metadataFilepath)
		err := os.MkdirAll(filepath.Dir(metadataAbsoluteFilepath), os.FileMode(0755))
		assert.NoError(t, err)

		if path != a2 {
			err = os.WriteFile(metadataAbsoluteFilepath, []byte{}, os.FileMode(0600))
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

func TestExerciseDir(t *testing.T) {
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
		dir, err := ws.ExerciseDir(test.path)
		if !test.ok {
			assert.Error(t, err, test.path)
			continue
		}
		assert.Equal(t, filepath.Join(ws.Dir, "exercise"), dir, test.path)
	}
}
