package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const cfgHomeKey = "EXERCISM_CONFIG_HOME"

type CommandTest struct {
	App            *cobra.Command
	Cmd            *cobra.Command
	InitFn         func()
	TmpDir         string
	Args           []string
	OriginalValues struct {
		ConfigHome string
		Args       []string
	}
}

func (test *CommandTest) Setup(t *testing.T) {
	dir, err := ioutil.TempDir("", "command-test")
	assert.NoError(t, err)

	test.TmpDir = dir
	test.OriginalValues.ConfigHome = os.Getenv(cfgHomeKey)
	test.OriginalValues.Args = os.Args

	os.Setenv(cfgHomeKey, test.TmpDir)

	os.Args = test.Args

	test.Cmd.ResetFlags()
	test.InitFn()

	test.App = &cobra.Command{}
	test.App.AddCommand(test.Cmd)
}

func (test *CommandTest) Teardown(t *testing.T) {
	os.Setenv(cfgHomeKey, test.OriginalValues.ConfigHome)
	os.Args = test.OriginalValues.Args
	if err := os.RemoveAll(test.TmpDir); err != nil {
		t.Fatal(err)
	}
}
