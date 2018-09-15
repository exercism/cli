package workspace

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	dir, err := ioutil.TempDir("", "solution")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	s1 := &Metadata{
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

	s2, err := NewMetadata(dir)
	assert.NoError(t, err)
	assert.Nil(t, s2.SubmittedAt)
	assert.Equal(t, s1, s2)

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	s2.SubmittedAt = &ts

	err = s2.Write(dir)
	assert.NoError(t, err)

	s3, err := NewMetadata(dir)
	assert.NoError(t, err)
	assert.Equal(t, s2, s3)
}

func TestSuffix(t *testing.T) {
	testCases := []struct {
		metadata Metadata
		suffix   string
	}{
		{
			metadata: Metadata{
				Exercise: "bat",
				Dir:      "",
			},
			suffix: "",
		},
		{
			metadata: Metadata{
				Exercise: "bat",
				Dir:      "/path/to/bat",
			},
			suffix: "",
		},
		{
			metadata: Metadata{
				Exercise: "bat",
				Dir:      "/path/to/bat-2",
			},
			suffix: "2",
		},
		{
			metadata: Metadata{
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

func TestMetadataString(t *testing.T) {
	testCases := []struct {
		metadata Metadata
		desc     string
	}{
		{
			metadata: Metadata{
				Track:    "elixir",
				Exercise: "secret-handshake",
				Handle:   "",
				Dir:      "",
			},
			desc: "elixir/secret-handshake",
		},
		{
			metadata: Metadata{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
			},
			desc: "cpp/clock",
		},
		{
			metadata: Metadata{
				Track:       "cpp",
				Exercise:    "clock",
				Handle:      "alice",
				IsRequester: true,
				Dir:         "/path/to/clock-2",
			},
			desc: "cpp/clock (2)",
		},
		{
			metadata: Metadata{
				Track:       "fsharp",
				Exercise:    "hello-world",
				Handle:      "bob",
				IsRequester: false,
			},
			desc: "fsharp/hello-world by @bob",
		},
		{
			metadata: Metadata{
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
