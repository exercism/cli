package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExerciseMetadatas(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solutions")

	paths := []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "bravo"),
		filepath.Join(root, "charlie"),
	}
	sx, err := NewExerciseMetadatas(paths)
	assert.NoError(t, err)

	if assert.Equal(t, 3, len(sx)) {
		assert.Equal(t, "alpha", sx[0].ID)
		assert.Equal(t, "bravo", sx[1].ID)
		assert.Equal(t, "charlie", sx[2].ID)
	}

	paths = []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "delta"),
		filepath.Join(root, "bravo"),
	}
	_, err = NewExerciseMetadatas(paths)
	assert.Error(t, err)
}
