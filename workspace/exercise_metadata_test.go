package workspace

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExerciseMetadata(t *testing.T) {
	dir, err := ioutil.TempDir("", "solution")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	s1 := &ExerciseMetadata{
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

	s2, err := NewExerciseMetadata(dir)
	assert.NoError(t, err)
	assert.Nil(t, s2.SubmittedAt)
	assert.Equal(t, s1, s2)

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	s2.SubmittedAt = &ts

	err = s2.Write(dir)
	assert.NoError(t, err)

	s3, err := NewExerciseMetadata(dir)
	assert.NoError(t, err)
	assert.Equal(t, s2, s3)
}

func TestSuffix(t *testing.T) {
	testCases := []struct {
		metadata ExerciseMetadata
		suffix   string
	}{
		{
			metadata: ExerciseMetadata{
				Exercise: "bat",
				Dir:      "",
			},
			suffix: "",
		},
		{
			metadata: ExerciseMetadata{
				Exercise: "bat",
				Dir:      "/path/to/bat",
			},
			suffix: "",
		},
		{
			metadata: ExerciseMetadata{
				Exercise: "bat",
				Dir:      "/path/to/bat-2",
			},
			suffix: "2",
		},
		{
			metadata: ExerciseMetadata{
				Exercise: "bat",
				Dir:      "/path/to/bat-200",
			},
			suffix: "200",
		},
	}

	for _, tc := range testCases {
		testName := "Suffix of '" + tc.metadata.Dir + "' should be " + tc.suffix
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.suffix, tc.metadata.Suffix(), testName)
		})
	}
}

func TestExerciseMetadataString(t *testing.T) {
	testCases := []struct {
		metadata ExerciseMetadata
		desc     string
	}{
		{
			metadata: ExerciseMetadata{
				Track:    "elixir",
				Exercise: "secret-handshake",
				Handle:   "",
				Dir:      "",
			},
			desc: "elixir/secret-handshake",
		},
		{
			metadata: ExerciseMetadata{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
			},
			desc: "cpp/clock",
		},
		{
			metadata: ExerciseMetadata{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
				Dir:         "/path/to/clock-2",
			},
			desc: "cpp/clock (2)",
		},
		{
			metadata: ExerciseMetadata{
				Track:       "fsharp",
				Exercise:    "hello-world",
				Handle:      "bob",
				IsRequester: false,
			},
			desc: "fsharp/hello-world by @bob",
		},
		{
			metadata: ExerciseMetadata{
				Track:       "haskell",
				Exercise:    "allergies",
				Handle:      "charlie",
				IsRequester: false,
				Dir:         "/path/to/allergies-2",
			},
			desc: "haskell/allergies (2) by @charlie",
		},
	}

	for _, tc := range testCases {
		testName := "should stringify to '" + tc.desc + "'"
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.desc, tc.metadata.String())
		})
	}
}
