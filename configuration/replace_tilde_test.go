package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandsTildeInExercismDirectory(t *testing.T) {
	expandedDir := ReplaceTilde("~/exercism/directory")
	assert.NotContains(t, "~", expandedDir)
}
