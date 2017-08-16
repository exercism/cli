package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSolutions(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solutions")

	paths := []string{
		filepath.Join(root, "alpha"),
		filepath.Join(root, "bravo"),
		filepath.Join(root, "charlie"),
	}
	sx, err := NewSolutions(paths)
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
	_, err = NewSolutions(paths)
	assert.Error(t, err)
}

func TestSolutions(t *testing.T) {
	solutions := []*Solution{
		{Track: "a", Exercise: "foo"},
		{Track: "b", Exercise: "foo"},
		{Track: "c", Exercise: "foo"},
	}
	sx := Solutions(solutions)
	display := "  [1] a/foo\n  [2] b/foo\n  [3] c/foo\n"
	assert.Equal(t, display, sx.Display())

	_, err := sx.Get(0)
	assert.Error(t, err)

	a, err := sx.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, "a", a.Track)

	b, err := sx.Get(2)
	assert.NoError(t, err)
	assert.Equal(t, "b", b.Track)

	c, err := sx.Get(3)
	assert.NoError(t, err)
	assert.Equal(t, "c", c.Track)

	_, err = sx.Get(4)
	assert.Error(t, err)
}
