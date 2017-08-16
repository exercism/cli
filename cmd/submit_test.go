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
	"github.com/stretchr/testify/assert"
)

func TestSubmit(t *testing.T) {
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

	cmdTest := &CommandTest{
		Cmd:    submitCmd,
		InitFn: initSubmitCmd,
		Args:   []string{"fakeapp", "submit", "bogus-exercise"},
	}
	cmdTest.Setup(t)
	defer cmdTest.Teardown(t)

	// Create a temp dir for the config and the exercise files.
	dir := filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))

	solution := &workspace.Solution{
		ID:          "bogus-solution-uuid",
		Track:       "bogus-track",
		Exercise:    "bogus-exercise",
		IsRequester: true,
	}
	err := solution.Write(dir)

	for _, file := range []file{file1, file2, file3} {
		err := ioutil.WriteFile(filepath.Join(dir, file.relativePath), []byte(file.contents), os.FileMode(0755))
		assert.NoError(t, err)
	}

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
			defer file.Close()

			if err != nil {
				t.Fatal(err)
			}
			body, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}
			submittedFiles[fileHeader.Filename] = string(body)
		}
	})
	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	// Create a fake user config.
	usrCfg := config.NewEmptyUserConfig()
	usrCfg.Workspace = cmdTest.TmpDir
	err = usrCfg.Write()
	assert.NoError(t, err)

	// Create a fake CLI config.
	cliCfg, err := config.NewCLIConfig()
	assert.NoError(t, err)
	cliCfg.Tracks["bogus-track"] = config.NewTrack("bogus-track")
	err = cliCfg.Write()
	assert.NoError(t, err)

	// Create a fake API config.
	apiCfg, err := config.NewAPIConfig()
	assert.NoError(t, err)
	apiCfg.BaseURL = ts.URL
	apiCfg.Endpoints["submit"] = "?%s"
	err = apiCfg.Write()
	assert.NoError(t, err)

	// Execute the command!
	cmdTest.App.Execute()

	// We got only the file we expected.
	assert.Equal(t, 2, len(submittedFiles))
	for _, file := range []file{file1, file2} {
		path := string(os.PathSeparator) + file.relativePath
		assert.Equal(t, file.contents, submittedFiles[path])
	}
}
