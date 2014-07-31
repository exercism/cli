package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandsTildeInExercismDirectory(t *testing.T) {
	expandedDir := ReplaceTilde("~/exercism/directory")
	assert.NotContains(t, "~", expandedDir)
}
