package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedact(t *testing.T) {
	fakeToken := "1a11111aaaa111aa1a11111a11111aa1"
	expected := "1a11*************************aa1"

	assert.Equal(t, expected, redact(fakeToken))
}
