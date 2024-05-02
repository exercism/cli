//go:build !windows

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultWorkspaceDir(t *testing.T) {
	testCases := []struct {
		cfg      Config
		expected string
	}{
		{
			cfg:      Config{OS: "darwin", Home: "/User/charlie", DefaultDirName: "apple"},
			expected: "/User/charlie/Apple",
		},
		{
			cfg:      Config{OS: "linux", Home: "/home/bob", DefaultDirName: "banana"},
			expected: "/home/bob/banana",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, DefaultWorkspaceDir(tc.cfg), fmt.Sprintf("Operating System: %s", tc.cfg.OS))
	}
}
