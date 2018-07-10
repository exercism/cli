package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptFilename(t *testing.T) {

	testCases := []struct {
		desc      string
		filenames []string
		expected  bool
	}{

		{"allowed filename", []string{"beacon.ext", "falcon.zip"}, true},
		{"ignored filename", []string{"beacon|txt", "falcon.txt", "proof"}, false},
	}

	track := &Track{
		IgnorePatterns: []string{
			"con[|.]txt",
			"pro.f",
		},
	}

	for _, tc := range testCases {
		for _, filename := range tc.filenames {
			t.Run(fmt.Sprintf("%s %s", tc.desc, filename), func(t *testing.T) {
				got, err := track.AcceptFilename(filename)
				assert.NoError(t, err, fmt.Sprintf("%s %s", tc.desc, filename))
				assert.Equal(t, tc.expected, got, fmt.Sprintf("should return %t for %s,  but got %t", tc.expected, tc.desc, got))
			})
		}
	}
}
