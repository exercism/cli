package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	// Let's not actually print to standard out while testing.
	Out = ioutil.Discard

	cmdTest := &CommandTest{
		Cmd:    downloadCmd,
		InitFn: initDownloadCmd,
		Args:   []string{"fakeapp", "download", "bogus-exercise"},
	}
	cmdTest.Setup(t)
	defer cmdTest.Teardown(t)

	// Write a fake user config setting the workspace to the temp dir.
	userCfg := config.NewEmptyUserConfig()
	userCfg.Workspace = cmdTest.TmpDir
	err := userCfg.Write()
	assert.NoError(t, err)

	payloadBody := `
	{
		"solution": {
			"id": "bogus-id",
			"user": {
				"handle": "alice",
				"is_requester": true
			},
			"exercise": {
				"id": "bogus-exercise",
				"instructions_url": "http://example.com/bogus-exercise",
				"auto_approve": false,
				"track": {
					"id": "bogus-track",
					"language": "Bogus Language"
				}
			},
			"file_download_base_url": "%s",
			"files": [
			"%s",
			"%s",
			"%s"
			],
			"iteration": {
				"submitted_at": "2017-08-21t10:11:12.130z"
			}
		}
	}
	`

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	path1 := "file-1.txt"
	mux.HandleFunc("/"+path1, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 1")
	})

	path2 := "subdir/file-2.txt"
	mux.HandleFunc("/"+path2, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 2")
	})

	path3 := "file-3.txt"
	mux.HandleFunc("/"+path3, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	})

	payloadBody = fmt.Sprintf(payloadBody, server.URL+"/", path1, path2, path3)
	mux.HandleFunc("/solutions/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})

	// Write a fake api config setting the base url to the test server.
	apiCfg := config.NewEmptyAPIConfig()
	apiCfg.BaseURL = server.URL
	err = apiCfg.Write()
	assert.NoError(t, err)

	tests := []struct {
		path     string
		contents string
	}{
		{
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", "file-1.txt"),
			contents: "this is file 1",
		},
		{
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", "subdir", "file-2.txt"),
			contents: "this is file 2",
		},
		{
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", ".solution.json"),
			contents: `{"track":"bogus-track","exercise":"bogus-exercise","id":"bogus-id","url":"","handle":"alice","is_requester":true}`,
		},
	}

	cmdTest.App.Execute()

	for _, test := range tests {
		b, err := ioutil.ReadFile(test.path)
		assert.NoError(t, err)
		assert.Equal(t, test.contents, string(b))
	}

	// It doesn't write the empty file.
	path := filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", path3)
	_, err = os.Lstat(path)
	assert.True(t, os.IsNotExist(err))
}
