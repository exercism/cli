package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestExerciseWithMetadataDoesNotOverwrite(t *testing.T) {
	// * Test setup: an exercise directory with a metadata file
	// * Verify: it doesn't get overwritten

	tmpDir, err := ioutil.TempDir("", "has-metadata")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	// writeFakeMetadata
	metadata := &workspace.ExerciseMetadata{
		ID:          "bogus-solution-uuid",
		Track:       "bogus-track",
		Exercise:    "bogus-exercise",
		URL:         "http://example.com/bogus-url",
		IsRequester: true,
	}
	err = metadata.Write(dir)

	// get metadata modtime
	exercise := workspace.NewExerciseFromDir(dir)
	fileInfo, err := os.Lstat(exercise.MetadataFilepath())
	preDoctorMetadataModTime := fileInfo.ModTime()

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	// configure fixup flag
	args := []string{
		"--fixup",
	}
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupDoctorFlags(flags)
	err = flags.Parse(args)
	assert.NoError(t, err)

	err = runDoctor(cfg, flags)
	assert.NoError(t, err)

	fileInfo, err = os.Lstat(exercise.MetadataFilepath())
	postDoctorMetadataModTime := fileInfo.ModTime()
	assert.Equal(t, preDoctorMetadataModTime, postDoctorMetadataModTime)
}

func TestExerciseWithoutMetadataWritesMetadataWithoutTouchingExerciseFiles(t *testing.T) {
	// It should not overwrite any existing exercise/solution files.
	// * Test setup: a text file with different text than test server
	// * Verify: text file doesn't get overwritten, metadata file gets written

	// TODO: necessary to fake a test server?

	tmpDir, err := ioutil.TempDir("", "no-metadata")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	// TODO: extract from cmd/download?
	ts := fakeDownloadServer("true", "bogus-team")
	defer ts.Close()

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)
	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	// configure fixup flag
	args := []string{
		"--fixup",
	}
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupDoctorFlags(flags)
	err = flags.Parse(args)

	assert.NoError(t, err)
	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))
	exercise := workspace.NewExerciseFromDir(dir)

	testFilepath := filepath.Join(dir, "file.txt")
	err = ioutil.WriteFile(testFilepath, []byte("This is a file."), os.FileMode(0755))
	fileInfo, err := os.Lstat(testFilepath)
	preDoctorTestFileModTime := fileInfo.ModTime()

	ok, err := exercise.HasMetadata()
	assert.NoError(t, err)
	assert.False(t, ok)

	err = runDoctor(cfg, flags)
	assert.NoError(t, err)

	fileInfo, err = os.Lstat(testFilepath)
	postDoctorTestFileModTime := fileInfo.ModTime()
	assert.Equal(t, preDoctorTestFileModTime, postDoctorTestFileModTime)

	ok, err = exercise.HasMetadata()
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestMigratesAllExercisesWithOutput(t *testing.T) {
	// * Test setup: create two fake tracks with no metadata, one with one exercise, one with two exercises
	// * Verify: they all got migrated, and STDERR contains a report of what was done

}

func TestDryRun(t *testing.T) {
	// * Test setup: same as "all available", but passes --dry-run flag
	// * Verify: none got migrated, STDERR contains a report of what would have been done
}
