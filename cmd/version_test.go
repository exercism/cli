package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exercism/cli/cli"
	"github.com/stretchr/testify/assert"
)

func TestCurrentVersion(t *testing.T) {
	expected := fmt.Sprintf("exercism version %s", Version)

	actual := currentVersion()
	assert.Equal(t, expected, actual)
}

func TestVersionUpdateCheck(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"tag_name": "v2.0.0"}`)
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()
	cli.LatestReleaseURL = ts.URL

	tests := []struct {
		version  string
		expected string
	}{
		{
			// It returns new version available for versions older than latest.
			version:  "1.0.0",
			expected: "A new CLI version is available. Run `exercism upgrade` to update to 2.0.0",
		},
		{
			// It returns up to date for versions matching latest.
			version:  "2.0.0",
			expected: "Your CLI version is up to date.",
		},
		{
			// It returns up to date for versions newer than latest.
			version:  "2.0.1",
			expected: "Your CLI version is up to date.",
		},
	}

	for _, test := range tests {
		c := &cli.CLI{
			Version: test.version,
		}

		actual, err := checkForUpdate(c)

		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Equal(t, test.expected, actual)
	}
}
