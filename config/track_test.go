package config

import (
	"fmt"
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

	for name, ok := range testCases {
		testName := fmt.Sprintf("%s is %s", name, acceptability(ok))
		t.Run(testName, func(t *testing.T) {
			acceptable, err := track.AcceptFilename(name)
			assert.NoError(t, err, name)
			assert.Equal(t, ok, acceptable, testName)
		})
	}
}

func acceptability(ok bool) string {
	if ok {
		return "fine"
	}
	return "not acceptable"
}
