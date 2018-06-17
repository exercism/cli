package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackIgnoreString(t *testing.T) {
	track := &Track{
		IgnorePatterns: []string{
			"con[.]txt",
			"pro.f",
		},
	}

	testCases := map[string]bool{
		"falcon.txt": false,
		"beacon|txt": true,
		"beacon.ext": true,
		"proof":      false,
	}

	for name, acceptable := range testCases {
		testName := name + " should " + notIfNeeded(acceptable) + "be an acceptable name."
		t.Run(testName, func(t *testing.T) {
			ok, err := track.AcceptFilename(name)
			assert.NoError(t, err, name)
			assert.Equal(t, acceptable, ok, testName)
		})
	}
}

func notIfNeeded(b bool) string {
	if !b {
		return "not "
	}
	return ""
}
