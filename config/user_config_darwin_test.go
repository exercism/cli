package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeWorkspace(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	cfg := &UserConfig{Home: "/home/alice"}
	tests := []struct {
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

	for _, test := range tests {
		cfg.Workspace = test.in
		cfg.Normalize()
		assert.Equal(t, test.out, cfg.Workspace)
	}
}
