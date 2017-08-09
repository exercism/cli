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

	s1 := &Solution{
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

func TestSuffix(t *testing.T) {
	tests := []struct {
		solution Solution
		suffix   string
	}{
		{
			solution: Solution{
				Exercise: "bat",
				Dir:      "",
			},
			suffix: "",
		},
		{
			solution: Solution{
				Exercise: "bat",
				Dir:      "/path/to/bat",
			},
			suffix: "",
		},
		{
			solution: Solution{
				Exercise: "bat",
				Dir:      "/path/to/bat-2",
			},
			suffix: "2",
		},
		{
			solution: Solution{
				Exercise: "bat",
				Dir:      "/path/to/bat-200",
			},
			suffix: "200",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.suffix, test.solution.Suffix())
	}
}

func TestSolutionString(t *testing.T) {
	tests := []struct {
		solution Solution
		desc     string
	}{
		{
			solution: Solution{
				Track:    "elixir",
				Exercise: "secret-handshake",
				Handle:   "",
				Dir:      "",
			},
			desc: "elixir/secret-handshake",
		},
		{
			solution: Solution{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
			},
			desc: "cpp/clock",
		},
		{
			solution: Solution{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
				Dir:         "/path/to/clock-2",
			},
			desc: "cpp/clock (2)",
		},
		{
			solution: Solution{
				Track:       "fsharp",
				Exercise:    "hello-world",
				Handle:      "bob",
				IsRequester: false,
			},
			desc: "fsharp/hello-world by @bob",
		},
		{
			solution: Solution{
				Track:       "haskell",
				Exercise:    "allergies",
				Handle:      "charlie",
				IsRequester: false,
				Dir:         "/path/to/allergies-2",
			},
			desc: "haskell/allergies (2) by @charlie",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.desc, test.solution.String())
	}
}
