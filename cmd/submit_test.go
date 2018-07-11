package cmd

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSubmitFiles(t *testing.T) {
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

	// Set up the test server.
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(2 << 10)
		if err != nil {
			t.Fatal(err)
		}
		mf := r.MultipartForm

		files := mf.File["files[]"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			body, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}
			submittedFiles[fileHeader.Filename] = string(body)
		}
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	tmpDir, err := ioutil.TempDir("", "submit-files")
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))

	type file struct {
		relativePath string
		contents     string
	}

	file1 := file{
		relativePath: "file-1.txt",
		contents:     "This is file 1.",
	}
	file2 := file{
		relativePath: filepath.Join("subdir", "file-2.txt"),
		contents:     "This is file 2.",
	}
	file3 := file{
		relativePath: "README.md",
		contents:     "The readme.",
	}

	filenames := make([]string, 0, 3)
	for _, file := range []file{file1, file2, file3} {
		path := filepath.Join(dir, file.relativePath)
		filenames = append(filenames, path)
		err := ioutil.WriteFile(path, []byte(file.contents), os.FileMode(0755))
		assert.NoError(t, err)
	}

	solution := &workspace.Solution{
		ID:          "bogus-solution-uuid",
		Track:       "bogus-track",
		Exercise:    "bogus-exercise",
		IsRequester: true,
	}
	err = solution.Write(dir)
	assert.NoError(t, err)

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupSubmitFlags(flags)
	flagArgs := []string{
		"--files",
		strings.Join(filenames, ","),
	}
	err = flags.Parse(flagArgs)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cliCfg := &config.CLIConfig{
		Config: config.New(tmpDir, "cli"),
		Tracks: config.Tracks{},
	}
	cliCfg.Tracks["bogus-track"] = config.NewTrack("bogus-track")
	err = cliCfg.Write()
	assert.NoError(t, err)

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
		CLIConfig:       cliCfg,
	}

	err = runSubmit(cfg, flags, []string{})
	assert.NoError(t, err)

	// We currently have a bug, and we're not filtering anything.
	// Fix that in a separate commit.
	assert.Equal(t, 3, len(submittedFiles))

	for _, file := range []file{file1, file2, file3} {
		path := string(os.PathSeparator) + file.relativePath
		assert.Equal(t, file.contents, submittedFiles[path])
	}
}
