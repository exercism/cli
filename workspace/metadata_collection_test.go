package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetadataCollection(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solutions")

	paths := []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "bravo"),
		filepath.Join(root, "charlie"),
	}
	collection, err := NewMetadataCollection(paths)
	assert.NoError(t, err)

	if assert.Equal(t, 3, len(collection)) {
		assert.Equal(t, "alpha", collection[0].ID)
		assert.Equal(t, "bravo", collection[1].ID)
		assert.Equal(t, "charlie", collection[2].ID)
	}

	paths = []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "delta"),
		filepath.Join(root, "bravo"),
	}
	_, err = NewMetadataCollection(paths)
	assert.Error(t, err)
}
