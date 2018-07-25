package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/exercism/cli/config"
	ws "github.com/exercism/cli/workspace"
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
		// It uses the default base API url to infer the host
		assert.Regexp(t, "exercism.io/my/settings", err.Error())
	}
}

func TestDownloadWithoutWorkspace(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")
	cfg := config.Configuration{
		UserViperConfig: v,
	}

	err := runDownload(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "re-run the configure", err.Error())
	}
}

func TestDownloadWithoutBaseURL(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", "/home/whatever")
	cfg := config.Configuration{
		UserViperConfig: v,
	}

	err := runDownload(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "re-run the configure", err.Error())
	}
}

func TestDownloadWithoutFlags(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", "/home/username")
	v.Set("apibaseurl", "http://example.com")

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

	testCases := []struct {
		requestor   string
		expectedDir string
		flags       map[string]string
	}{
		{
			requestor:   requestorSelf,
			expectedDir: "",
			flags:       map[string]string{"exercise": "bogus-exercise"},
		},
		{
			requestor:   requestorSelf,
			expectedDir: "",
			flags:       map[string]string{"uuid": "bogus-id"},
		},
		{
			requestor:   requestorOther,
			expectedDir: filepath.Join("users", "alice"),
			flags:       map[string]string{"uuid": "bogus-id"},
		},
	}

	for _, tc := range testCases {
		tmpDir, err := ioutil.TempDir("", "download-cmd")
		defer os.RemoveAll(tmpDir)
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
		for name, value := range tc.flags {
			flags.Set(name, value)
		}

		err = runDownload(cfg, flags, []string{})
		assert.NoError(t, err)

		targetDir := filepath.Join(tmpDir, tc.expectedDir)
		assertDownloadedCorrectFiles(t, targetDir, tc.requestor)

		metadata := `{
			"track": "bogus-track",
			"exercise":"bogus-exercise",
			"id":"bogus-id",
			"url":"",
			"handle":"alice",
			"is_requester":%s,
			"auto_approve":false
		}`
		metadata = fmt.Sprintf(metadata, tc.requestor)
		metadata = compact(t, metadata)

		path := filepath.Join(targetDir, "bogus-track", "bogus-exercise", ws.SolutionMetadataFilepath())
		b, err := ioutil.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, metadata, string(b), "the solution metadata file")
	}
}

func fakeDownloadServer(requestor string) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/file-1.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 1")
	})

	mux.HandleFunc("/subdir/file-2.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 2")
	})

	mux.HandleFunc("/file-3.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	})

	payloadBody := fmt.Sprintf(payloadTemplate, requestor, server.URL+"/")
	mux.HandleFunc("/solutions/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})
	mux.HandleFunc("/solutions/bogus-id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, payloadBody)
	})

	return server
}

func assertDownloadedCorrectFiles(t *testing.T, targetDir, requestor string) {
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
		"/file-1.txt",
		"/subdir/file-2.txt",
		"/file-3.txt"
		],
		"iteration": {
			"submitted_at": "2017-08-21t10:11:12.130z"
		}
	}
}
`

func compact(t *testing.T, s string) string {
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(s))
	assert.NoError(t, err)
	return buffer.String()
}
