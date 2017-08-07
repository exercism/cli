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
		SubmittedAt: time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC),
	}
	err = s1.Write(dir)
	assert.NoError(t, err)

	s2, err := NewSolution(dir)
	assert.NoError(t, err)

	assert.Equal(t, s1, s2)
}
