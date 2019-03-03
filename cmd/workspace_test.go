package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWorkspaceCmdWithoutWorkspace(t *testing.T) {
	v := viper.New()
	v.Set("Token", "abc123")
	err := runWorkspace(v, pflag.NewFlagSet("fake", pflag.PanicOnError))
	if assert.Error(t, err) {
		assert.Regexp(t, "re-run the configure", err.Error())
	}
}
func TestWorkspaceCmdWithNonexistentTeam(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fake-workspace")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://example.com")

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupWorkspaceFlags(flags)
	flags.Set("team", "NonExistentTeam")

	err = runWorkspace(v, flags)
	if assert.Error(t, err) {
		assert.Regexp(t, "team not found", err.Error())
	}
}

func TestWorkspaceCmd(t *testing.T) {
	co := newCapturedOutput()
	defer co.reset()

	tmpDir, err := ioutil.TempDir("", "fake-workspace")
	setupTrackandExercises(tmpDir, t)
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://example.com")

	testCases := []struct {
		desc     string
		args     []string
		expected string
	}{
		{
			desc:     "It prints the workspace path",
			args:     []string{},
			expected: tmpDir,
		},
		{
			desc:     "It prints the team workspace",
			args:     []string{"--team", "alice"},
			expected: filepath.Join(tmpDir, "teams/alice"),
		},
		{
			desc:     "It prints the tracks",
			args:     []string{"--track"},
			expected: filepath.Join(tmpDir, "track-a"),
		},
		{
			desc:     "It prints the team's tracks",
			args:     []string{"--track", "--team", "alice"},
			expected: filepath.Join(tmpDir, "teams/alice/track-a"),
		},
		{
			desc:     "It prints the exercises",
			args:     []string{"--exercise"},
			expected: filepath.Join(tmpDir, "track-a", "exercise-one"),
		},
		{
			desc:     "It prints the team's exercises",
			args:     []string{"--exercise", "--team", "alice"},
			expected: filepath.Join(tmpDir, "teams/alice/track-a/exercise-one"),
		},
	}

	for _, tc := range testCases {
		co.newOut = &bytes.Buffer{}
		co.override()

		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		setupWorkspaceFlags(flags)

		err := flags.Parse(tc.args)
		assert.NoError(t, err)

		err = runWorkspace(v, flags)
		assert.NoError(t, err)

		assert.Regexp(t, tc.expected, Out, tc.desc)
	}
}

func setupTrackandExercises(tmpDir string, t *testing.T) {
	a1 := filepath.Join(tmpDir, "track-a", "exercise-one")

	teamDir := filepath.Join(tmpDir, "teams", "alice", "track-a", "exercise-one")

	for _, path := range []string{a1, teamDir} {
		metadataAbsoluteFilepath := filepath.Join(path, ".exercism/metadata.json")
		err := os.MkdirAll(filepath.Dir(metadataAbsoluteFilepath), os.FileMode(0755))
		assert.NoError(t, err)

		err = ioutil.WriteFile(metadataAbsoluteFilepath, []byte{}, os.FileMode(0600))
		assert.NoError(t, err)
	}
}
