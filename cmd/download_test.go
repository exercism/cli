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

const payloadTemplate = `
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

func TestDownload(t *testing.T) {
	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	Err = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	cmdTest := &CommandTest{
		Cmd:    downloadCmd,
		InitFn: initDownloadCmd,
		Args:   []string{"fakeapp", "download", "--exercise=bogus-exercise"},
	}
	cmdTest.Setup(t)
	defer cmdTest.Teardown(t)

	mockServer := makeMockServer()
	defer mockServer.Close()

	err := writeFakeUserConfigSettings(cmdTest.TmpDir, mockServer.URL)
	assert.NoError(t, err)

	testCases := []struct {
		desc     string
		path     string
		contents string
	}{
		{
			desc:     "It should download a file.",
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", "file-1.txt"),
			contents: "this is file 1",
		},
		{
			desc:     "It should download a file in a subdir.",
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", "subdir", "file-2.txt"),
			contents: "this is file 2",
		},
		{
			desc:     "It creates the .solution.json file.",
			path:     filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", ".solution.json"),
			contents: `{"track":"bogus-track","exercise":"bogus-exercise","id":"bogus-id","url":"","handle":"alice","is_requester":true,"auto_approve":false}`,
		},
	}

	cmdTest.App.Execute()

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			b, err := ioutil.ReadFile(tc.path)
			assert.NoError(t, err)
			assert.Equal(t, tc.contents, string(b))
		})
	}

	path := filepath.Join(cmdTest.TmpDir, "bogus-track", "bogus-exercise", "file-3.txt")
	_, err = os.Lstat(path)
	assert.True(t, os.IsNotExist(err), "It should not write the file if empty.")
}

func writeFakeUserConfigSettings(tmpDirPath, serverURL string) error {
	userCfg := config.NewEmptyUserConfig()
	userCfg.Workspace = tmpDirPath
	userCfg.APIBaseURL = serverURL
	return userCfg.Write()
}

func makeMockServer() *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

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

	payloadBody := fmt.Sprintf(payloadTemplate, server.URL+"/", path1, path2, path3)
	mux.HandleFunc("/solutions/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})

	return server

}
