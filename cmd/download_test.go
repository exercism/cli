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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDownloadWithoutToken(t *testing.T) {
	cfg := config.Configuration{
		UserViperConfig: viper.New(),
	}

	err := runDownload(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "Welcome to Exercism", err.Error())
	}
}

func TestDownloadWithoutFlags(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")

	cfg := config.Configuration{
		UserViperConfig: v,
	}

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupDownloadFlags(flags)

	err := runDownload(cfg, flags, []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "need an --exercise name or a solution --uuid", err.Error())
	}
}

func TestDownloadTHISONE(t *testing.T) {
	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	Err = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	testCases := []struct {
		requestor       string
		expectedDir     string
		flag, flagValue string
	}{
		{requestorSelf, "", "exercise", "bogus-exercise"},
		{requestorSelf, "", "uuid", "bogus-id"},
		{requestorOther, filepath.Join("users", "alice"), "uuid", "bogus-id"},
	}

	for _, tc := range testCases {
		tmpDir, err := ioutil.TempDir("", "download-cmd")
		assert.NoError(t, err)

		ts := fakeDownloadServer(tc.requestor)
		defer ts.Close()

		v := viper.New()
		v.Set("workspace", tmpDir)
		v.Set("apibaseurl", ts.URL)
		v.Set("token", "abc123")

		cfg := config.Configuration{
			UserViperConfig: v,
		}
		flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
		setupDownloadFlags(flags)
		flags.Set(tc.flag, tc.flagValue)

		err = runDownload(cfg, flags, []string{})
		assert.NoError(t, err)

		assertDownloadedCorrectFiles(t, filepath.Join(tmpDir, tc.expectedDir), tc.requestor)
	}
}

func fakeDownloadServer(requestor string) *httptest.Server {
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

	payloadBody := fmt.Sprintf(payloadTemplate, requestor, server.URL+"/", path1, path2, path3)
	mux.HandleFunc("/solutions/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})
	mux.HandleFunc("/solutions/bogus-id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})

	return server
}

func assertDownloadedCorrectFiles(t *testing.T, targetDir, requestor string) {
	metadata := `{"track":"bogus-track","exercise":"bogus-exercise","id":"bogus-id","url":"","handle":"alice","is_requester":%s,"auto_approve":false}`
	expectedFiles := []struct {
		desc     string
		path     string
		contents string
	}{
		{
			desc:     "a file in the exercise root directory",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "file-1.txt"),
			contents: "this is file 1",
		},
		{
			desc:     "a file in a subdirectory",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "subdir", "file-2.txt"),
			contents: "this is file 2",
		},
		{
			desc:     "the solution metadata file",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", ".solution.json"),
			contents: fmt.Sprintf(metadata, requestor),
		},
	}

	for _, file := range expectedFiles {
		t.Run(file.desc, func(t *testing.T) {
			b, err := ioutil.ReadFile(file.path)
			assert.NoError(t, err)
			assert.Equal(t, file.contents, string(b))
		})
	}

	path := filepath.Join(targetDir, "bogus-track", "bogus-exercise", "file-3.txt")
	_, err := os.Lstat(path)
	assert.True(t, os.IsNotExist(err), "It should not write the file if empty.")
}

const requestorSelf = "true"
const requestorOther = "false"

const payloadTemplate = `
{
	"solution": {
		"id": "bogus-id",
		"user": {
			"handle": "alice",
			"is_requester": %s
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
