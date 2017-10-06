package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUpToDate(t *testing.T) {
	tests := []struct {
		cliVersion string
		releaseTag string
		ok         bool
	}{
		{
			// It returns false for versions less than release.
			cliVersion: "1.0.0",
			releaseTag: "v1.0.1",
			ok:         false,
		},
		{
			// It returns false for pre-release versions of release.
			cliVersion: "1.0.1-alpha.1",
			releaseTag: "v1.0.1",
			ok:         false,
		},
		{
			// It returns true for versions equal to release.
			cliVersion: "2.0.1",
			releaseTag: "v2.0.1",
			ok:         true,
		},
		{
			// It returns true for versions greater than release.
			cliVersion: "2.0.2",
			releaseTag: "v2.0.1",
			ok:         true,
		},
	}

	for _, test := range tests {
		c := &CLI{
			Version:       test.cliVersion,
			LatestRelease: &Release{TagName: test.releaseTag},
		}

		ok, err := c.IsUpToDate()
		assert.NoError(t, err)
		assert.Equal(t, test.ok, ok, test.cliVersion)
	}
}

func TestIsUpToDateWithoutRelease(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"tag_name": "v2.0.0"}`)
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()
	LatestReleaseURL = ts.URL

	c := &CLI{
		Version: "1.0.0",
	}

	ok, err := c.IsUpToDate()
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.NotNil(t, c.LatestRelease)
}
