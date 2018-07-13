package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSubmitRelativePath(t *testing.T) {
	//t.Skip("The Windows build is failing and needs to be debugged.\nSee https://ci.appveyor.com/project/kytrinyx/cli/build/110")

	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	Err = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()
	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := ioutil.TempDir("", "relative-path")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeSolution(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	err = ioutil.WriteFile(filepath.Join(dir, "file.txt"), []byte("This is a file."), os.FileMode(0755))

	err = os.Chdir(dir)
	assert.NoError(t, err)

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{"file.txt"})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(submittedFiles))
	assert.Equal(t, "This is a file.", submittedFiles["\\file.txt"])
}
