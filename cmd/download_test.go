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

func TestDownload(t *testing.T) {
	oldOut := Out
	oldErr := Err
	Out = ioutil.Discard
	Err = ioutil.Discard
	defer func() {
		Out = oldOut
		Err = oldErr
	}()

	tmpDir, err := ioutil.TempDir("", "download-cmd")
	assert.NoError(t, err)

	ts := fakeDownloadServer()
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
	flags.Set("exercise", "bogus-exercise")

	err = runDownload(cfg, flags, []string{})
	assert.NoError(t, err)

	expectedFiles := []struct {
		desc     string
		path     string
		contents string
	}{
		{
			desc:     "a file in the exercise root directory",
			path:     filepath.Join(tmpDir, "bogus-track", "bogus-exercise", "file-1.txt"),
			contents: "this is file 1",
		},
		{
			desc:     "a file in a subdirectory",
			path:     filepath.Join(tmpDir, "bogus-track", "bogus-exercise", "subdir", "file-2.txt"),
			contents: "this is file 2",
		},
		{
			desc:     "the solution metadata file",
			path:     filepath.Join(tmpDir, "bogus-track", "bogus-exercise", ".solution.json"),
			contents: `{"track":"bogus-track","exercise":"bogus-exercise","id":"bogus-id","url":"","handle":"alice","is_requester":true,"auto_approve":false}`,
		},
	}

	for _, file := range expectedFiles {
		t.Run(file.desc, func(t *testing.T) {
			b, err := ioutil.ReadFile(file.path)
			assert.NoError(t, err)
			assert.Equal(t, file.contents, string(b))
		})
	}

	path := filepath.Join(tmpDir, "bogus-track", "bogus-exercise", "file-3.txt")
	_, err = os.Lstat(path)
	assert.True(t, os.IsNotExist(err), "It should not write the file if empty.")
}

func fakeDownloadServer() *httptest.Server {
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
