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
	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: viper.New(),
		DefaultBaseURL:  "http://example.com",
	}

	err := runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	assert.Regexp(t, "Welcome to Exercism", err.Error())
}

func TestSubmitWithoutWorkspace(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")

	cfg := config.Configuration{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err := runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	assert.Regexp(t, "run configure", err.Error())
}

func TestSubmitNonExistentFile(t *testing.T) {
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
	files := []string{
		filepath.Join(tmpDir, "file-1.txt"),
		"no-such-file.txt",
		filepath.Join(tmpDir, "file-2.txt"),
	}
	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	assert.Regexp(t, "no such file", err.Error())
}

func TestSubmitFilesAndDir(t *testing.T) {
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
	files := []string{
		filepath.Join(tmpDir, "file-1.txt"),
		tmpDir,
		filepath.Join(tmpDir, "file-2.txt"),
	}
	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
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
	writeFakeSolution(t, dir, "bogus-track", "bogus-exercise")

	file1 := filepath.Join(dir, "file-1.txt")
	err = ioutil.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	file2 := filepath.Join(dir, "subdir", "file-2.txt")
	err = ioutil.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))
	assert.NoError(t, err)

	// We don't filter *.md files if you explicitly pass the file path.
	readme := filepath.Join(dir, "README.md")
	err = ioutil.WriteFile(readme, []byte("This is the readme."), os.FileMode(0755))
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

	files := []string{
		file1, file2, readme,
	}
	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(submittedFiles))

	assert.Equal(t, "This is file 1.", submittedFiles[string(os.PathSeparator)+"file-1.txt"])
	assert.Equal(t, "This is file 2.", submittedFiles[string(os.PathSeparator)+filepath.Join("subdir", "file-2.txt")])
	assert.Equal(t, "This is the readme.", submittedFiles[string(os.PathSeparator)+"README.md"])
}

func TestSubmitFilesFromDifferentSolutions(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "dir-1-submit")
	assert.NoError(t, err)

	dir1 := filepath.Join(tmpDir, "bogus-track", "bogus-exercise-1")
	os.MkdirAll(dir1, os.FileMode(0755))
	writeFakeSolution(t, dir1, "bogus-track", "bogus-exercise-1")

	dir2 := filepath.Join(tmpDir, "bogus-track", "bogus-exercise-2")
	os.MkdirAll(dir2, os.FileMode(0755))
	writeFakeSolution(t, dir2, "bogus-track", "bogus-exercise-2")

	file1 := filepath.Join(dir1, "file-1.txt")
	err = ioutil.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	file2 := filepath.Join(dir2, "file-2.txt")
	err = ioutil.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)

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

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1, file2})
	assert.Error(t, err)
	assert.Regexp(t, "more than one solution", err.Error())
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
