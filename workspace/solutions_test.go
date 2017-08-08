package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolutions(t *testing.T) {
	solutions := []Solution{
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
