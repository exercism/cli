package cmd

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSubmitWithoutToken(t *testing.T) {
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: viper.New(),
		DefaultBaseURL:  "http://example.com",
	}

	err := runSubmit(cfg, flags, []string{})
	assert.Regexp(t, "Welcome to Exercism", err.Error())
}

func TestSubmitWithoutWorkspace(t *testing.T) {
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)

	v := viper.New()
	v.Set("token", "abc123")

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err := runSubmit(cfg, flags, []string{})
	assert.Regexp(t, "run configure", err.Error())
}

func TestSubmitNonExistentFile(t *testing.T) {
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)

	tmpDir, err := ioutil.TempDir("", "submit-no-such-file")
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err = ioutil.WriteFile(filepath.Join(tmpDir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tmpDir, "file-2.txt"), []byte("This is file 2"), os.FileMode(0755))
	assert.NoError(t, err)

	err = runSubmit(cfg, flags, []string{filepath.Join(tmpDir, "file-1.txt"), "no-such-file.txt", filepath.Join(tmpDir, "file-2.txt")})
	assert.Regexp(t, "no such file", err.Error())
}

func TestSubmitFilesAndDir(t *testing.T) {
	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)

	tmpDir, err := ioutil.TempDir("", "submit-no-such-file")
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err = ioutil.WriteFile(filepath.Join(tmpDir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(tmpDir, "file-2.txt"), []byte("This is file 2"), os.FileMode(0755))
	assert.NoError(t, err)

	err = runSubmit(cfg, flags, []string{filepath.Join(tmpDir, "file-1.txt"), tmpDir, filepath.Join(tmpDir, "file-2.txt")})
	assert.Regexp(t, "is a directory", err.Error())
}

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
	ts := fakeSubmitServer(t, submittedFiles)
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

	writeFakeSolution(t, dir, "bogus-track", "bogus-exercise")

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupSubmitFlags(flags)

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

	err = runSubmit(cfg, flags, filenames)
	assert.NoError(t, err)

	// We currently have a bug, and we're not filtering anything.
	// Fix that in a separate commit.
	assert.Equal(t, 3, len(submittedFiles))

	for _, file := range []file{file1, file2, file3} {
		path := string(os.PathSeparator) + file.relativePath
		assert.Equal(t, file.contents, submittedFiles[path])
	}
}

func fakeSubmitServer(t *testing.T, submittedFiles map[string]string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	return httptest.NewServer(handler)
}

func writeFakeSolution(t *testing.T, dir, trackID, exerciseSlug string) {
	solution := &workspace.Solution{
		ID:          "bogus-solution-uuid",
		Track:       trackID,
		Exercise:    exerciseSlug,
		URL:         "http://example.com/bogus-url",
		IsRequester: true,
	}
	err := solution.Write(dir)
	assert.NoError(t, err)
}
