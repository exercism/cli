//go:build !windows

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
		{"/home/alice///foobar", "/home/alice/foobar"},
		{"~/foobar", "/home/alice/foobar"},
		{"/foobar/~/noexpand", "/foobar/~/noexpand"},
		{"/no/modification", "/no/modification"},
		{"relative", filepath.Join(cwd, "relative")},
		{"relative///path", filepath.Join(cwd, "relative", "path")},
	}

	for _, tc := range testCases {
		testName := "'" + tc.in + "' should be normalized as '" + tc.out + "'"
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.out, Resolve(tc.in, "/home/alice"), testName)
		})
	}
}
