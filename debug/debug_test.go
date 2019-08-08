package debug

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerboseEnabled(t *testing.T) {
	b := &bytes.Buffer{}
	output = b
	Verbose = true

	Println("World")
	if b.String() != "World\n" {
		t.Error("expected 'World' got", b.String())
	}
}

func TestVerboseDisabled(t *testing.T) {
	b := &bytes.Buffer{}
	output = b
	Verbose = false

	Println("World")
	if b.String() != "" {
		t.Error("expected '' got", b.String())
	}
}

func TestRedact(t *testing.T) {
	fakeToken := "1a11111aaaa111aa1a11111a11111aa1"
	expected := "1a11*************************aa1"

	assert.Equal(t, expected, Redact(fakeToken))
}
