package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExerciseMetadataCollection(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solutions")

	paths := []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "bravo"),
		filepath.Join(root, "charlie"),
	}
	metadata, err := NewExerciseMetadataCollection(paths)
	assert.NoError(t, err)

	if assert.Equal(t, 3, len(metadata)) {
		assert.Equal(t, "alpha", metadata[0].ID)
		assert.Equal(t, "bravo", metadata[1].ID)
		assert.Equal(t, "charlie", metadata[2].ID)
	}

	paths = []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "delta"),
		filepath.Join(root, "bravo"),
	}
	_, err = NewExerciseMetadataCollection(paths)
	assert.Error(t, err)
}
