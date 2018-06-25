package cmd

import (
	"testing"

	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestCensor(t *testing.T) {
	fakeToken := "1a11111aaaa111aa1a11111a11111aa1"
	expected := "1a11*************************aa1"

	c := cli.New("")

	cfg := config.NewEmptyUserConfig()
	cfg.Token = fakeToken

	status := NewStatus(c, *cfg)
	status.Censor = true
	status.Check()

	assert.Equal(t, expected, status.Configuration.Token)
}
