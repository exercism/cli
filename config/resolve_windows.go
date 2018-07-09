package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	testCases := []struct {
		in, out string
	}{
		{"", ""}, // don't make wild guesses
		{"C:\\alice\\\\foobar", "C:\\alice\\\\foobar"},
		{"\\foobar\\~\\noexpand", "\\foobar\\~\\noexpand"},
		{"\\no\\modification", "\\no\\modification"},
		{"relative", filepath.Join(cwd, "relative")},
		{"relative\\path", filepath.Join(cwd, "relative", "path")},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			desc := "'" + tc.in + "' should be normalized as '" + tc.out + "'"
			assert.Equal(t, tc.out, Resolve(tc.in, ""), desc)
		})
	}
}
