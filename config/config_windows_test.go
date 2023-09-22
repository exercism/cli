//go:build windows

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultWindowsWorkspaceDir(t *testing.T) {
	cfg := Config{OS: "windows", Home: "C:\\Something", DefaultDirName: "basename"}
	assert.Equal(t, "C:\\Something\\Basename", DefaultWorkspaceDir(cfg))
}
