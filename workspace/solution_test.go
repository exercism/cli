package workspace

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSolution(t *testing.T) {
	dir, err := ioutil.TempDir("", "solution")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	s1 := Solution{
		Track:       "a-track",
		Exercise:    "bogus-exercise",
		ID:          "abc",
		URL:         "http://example.com",
		Handle:      "alice",
		IsRequester: true,
		Dir:         dir,
	}
	err = s1.Write(dir)
	assert.NoError(t, err)

	s2, err := NewSolution(dir)
	assert.NoError(t, err)
	assert.Nil(t, s2.SubmittedAt)
	assert.Equal(t, s1, s2)

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	s2.SubmittedAt = &ts

	err = s2.Write(dir)
	assert.NoError(t, err)

	s3, err := NewSolution(dir)
	assert.NoError(t, err)
	assert.Equal(t, s2, s3)
}
