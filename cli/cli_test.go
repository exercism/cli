package cli

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUpgradeNeeded(t *testing.T) {
	tests := []struct {
		cliVersion string
		releaseTag string
		needed     bool
	}{
		{"1.0.0", "v1.0.1", true},
		{"2.0.1", "v2.0.1", false},
	}

	for _, test := range tests {
		c := &CLI{
			Version:       test.cliVersion,
			LatestRelease: &Release{TagName: test.releaseTag},
		}

		needed, err := c.IsUpgradeNeeded()
		assert.NoError(t, err)
		assert.Equal(t, test.needed, needed, test.cliVersion)
	}
}

func TestIsUpgradeNeededWithoutRelease(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"tag_name": "v2.0.0"}`)
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()
	LatestReleaseURL = ts.URL

	c := &CLI{
		Version: "1.0.0",
	}

	needed, err := c.IsUpgradeNeeded()
	assert.NoError(t, err)
	assert.True(t, needed)
	assert.NotNil(t, c.LatestRelease)
}
