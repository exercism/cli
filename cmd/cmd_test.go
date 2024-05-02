package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const cfgHomeKey = "EXERCISM_CONFIG_HOME"

// CommandTest makes it easier to write tests for Cobra commands.
//
// To initialize, give it the three fields Cmd, InitFn, and Args.
// Then call Setup, and defer the Teardown.
// The args are the faked out os.Args. The first two arguments
// in the Args will be ignored. These represent the command (e.g. exercism)
// and the subcommand (e.g. download).
// Pass any interactive responses needed for the test in a single
// String in MockInput, delimited by newlines.
//
// Finally, when you have done whatever other setup you need in your
// test, call the command by calling Execute on the App.
//
// Example:
//
//	cmdTest := &CommandTest{
//		Cmd:    myCmd,
//		InitFn: initMyCmd,
//		Args:   []string{"fakeapp", "mycommand", "arg1", "--flag", "value"},
//		MockInteractiveResponse: "first-input\nsecond\n",
//	}
//
// cmdTest.Setup(t)
// defer cmdTest.Teardown(t)
// ...
// cmdTest.App.Execute()
type CommandTest struct {
	App                     *cobra.Command
	Cmd                     *cobra.Command
	InitFn                  func()
	TmpDir                  string
	Args                    []string
	MockInteractiveResponse string
	OriginalValues          struct {
		ConfigHome string
		Args       []string
	}
}

// Setup does all the prep and initialization for testing a command.
// It creates a fake Cobra app to provide a clean harness for the test,
// and adds the command under test to it as a subcommand.
// It also resets and reconfigures the command under test to
// make sure we're not getting any accidental pollution from the existing
// environment or other tests. Lastly, because we need to override some of
// the global environment settings, the setup method also stores the existing
// values so that Teardown can set them back the way they were when the test
// has completed.
// The method takes a *testing.T as an argument, that way the method can
// fail the test if the creation of the temporary directory fails.
func (test *CommandTest) Setup(t *testing.T) {
	dir, err := os.MkdirTemp("", "command-test")
	defer os.RemoveAll(dir)
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
	test.App.SetOutput(Err)
}

// Teardown puts the environment back the way it was before the test.
// The method takes a *testing.T so that it can blow up if it fails to
// clean up after itself.
func (test *CommandTest) Teardown(t *testing.T) {
	os.Setenv(cfgHomeKey, test.OriginalValues.ConfigHome)
	os.Args = test.OriginalValues.Args
	if err := os.RemoveAll(test.TmpDir); err != nil {
		t.Fatal(err)
	}
}

// capturedOutput lets us more easily redirect streams in the tests.
type capturedOutput struct {
	oldOut, oldErr, newOut, newErr io.Writer
}

// newCapturedOutput creates a new value to override the streams.
func newCapturedOutput() capturedOutput {
	return capturedOutput{
		oldOut: Out,
		oldErr: Err,
		newOut: io.Discard,
		newErr: io.Discard,
	}
}

// override sets the package variables to the fake streams.
func (co capturedOutput) override() {
	Out = co.newOut
	Err = co.newErr
}

// reset puts back the original streams for the commands to write to.
func (co capturedOutput) reset() {
	Out = co.oldOut
	Err = co.oldErr
}
