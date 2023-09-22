//go:build !windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSubmitFilesInSymlinkedPath(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "symlink-destination")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)
	dstDir := filepath.Join(tmpDir, "workspace")

	srcDir, err := os.MkdirTemp("", "symlink-source")
	defer os.RemoveAll(srcDir)
	assert.NoError(t, err)

	err = os.Symlink(srcDir, dstDir)
	assert.NoError(t, err)

	dir := filepath.Join(dstDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", dstDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(filepath.Join(dir, "file.txt"), []byte("This is a file."), os.FileMode(0755))
	assert.NoError(t, err)

	err = runSubmit(cfg, pflag.NewFlagSet("symlinks", pflag.PanicOnError), []string{file})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(submittedFiles))
	assert.Equal(t, "This is a file.", submittedFiles["file.txt"])
}
