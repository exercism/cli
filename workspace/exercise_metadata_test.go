package workspace

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	dir, err := ioutil.TempDir("", "metadata")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	metadata1 := &Metadata{
		Track:       "a-track",
		Exercise:    "bogus-exercise",
		ID:          "abc",
		URL:         "http://example.com",
		Handle:      "alice",
		IsRequester: true,
		Dir:         dir,
	}
	err = metadata1.Write(dir)
	assert.NoError(t, err)

	metadata2, err := NewMetadata(dir)
	assert.NoError(t, err)
	assert.Nil(t, metadata2.SubmittedAt)
	assert.Equal(t, metadata1, metadata2)

	ts := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	metadata2.SubmittedAt = &ts

	err = metadata2.Write(dir)
	assert.NoError(t, err)

	metadata3, err := NewMetadata(dir)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, metadata3)
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
