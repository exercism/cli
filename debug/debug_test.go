package debug

import (
	"bytes"
	"testing"
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
