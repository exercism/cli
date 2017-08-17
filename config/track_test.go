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

	tests := map[string]bool{
		"falcon.txt": false,
		"beacon|txt": true,
		"beacon.ext": true,
		"proof":      false,
	}

	for name, acceptable := range tests {
		ok, err := track.AcceptFilename(name)
		assert.NoError(t, err, name)
		assert.Equal(t, acceptable, ok, name)
	}
}
