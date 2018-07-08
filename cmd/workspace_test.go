package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWorkspace(t *testing.T) {
	Out = new(bytes.Buffer)

	v := viper.New()
	v.Set("workspace", "/home/alice")

	workspaceRun(v)

	b := Out.(*bytes.Buffer)
	assert.Equal(t, string("/home/alice\n"), b.String())
}
