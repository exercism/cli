package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDownloadWithoutToken(t *testing.T) {
	cfg := config.Config{
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
	cfg := config.Config{
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
	cfg := config.Config{
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

	cfg := config.Config{
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
		requester   bool
		expectedDir string
		flags       map[string]string
	}{
		{
			requester:   true,
			expectedDir: "",
			flags:       map[string]string{"exercise": "bogus-exercise"},
		},
		{
			requester:   true,
			expectedDir: "",
			flags:       map[string]string{"uuid": "bogus-id"},
		},
		{
			requester:   false,
			expectedDir: filepath.Join("users", "alice"),
			flags:       map[string]string{"uuid": "bogus-id"},
		},
		{
			requester:   true,
			expectedDir: filepath.Join("teams", "bogus-team"),
			flags:       map[string]string{"exercise": "bogus-exercise", "track": "bogus-track", "team": "bogus-team"},
		},
	}

	for _, tc := range testCases {
		tmpDir, err := ioutil.TempDir("", "download-cmd")
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		ts := fakeDownloadServer(strconv.FormatBool(tc.requester), tc.flags["team"])
		defer ts.Close()

		v := viper.New()
		v.Set("workspace", tmpDir)
		v.Set("apibaseurl", ts.URL)
		v.Set("token", "abc123")

		cfg := config.Config{
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
		assertDownloadedCorrectFiles(t, targetDir)

		dir := filepath.Join(targetDir, "bogus-track", "bogus-exercise")
		b, err := ioutil.ReadFile(workspace.NewExerciseFromDir(dir).MetadataFilepath())
		var s workspace.Metadata
		err = json.Unmarshal(b, &s)
		assert.NoError(t, err)

		assert.Equal(t, "bogus-track", s.Track)
		assert.Equal(t, "bogus-exercise", s.Exercise)
		assert.Equal(t, tc.requester, s.IsRequester)
	}
}

func fakeDownloadServer(requestor, teamSlug string) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/file-1.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 1")
	})

	mux.HandleFunc("/subdir/file-2.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is file 2")
	})

	mux.HandleFunc("/full/path/with/numeric-suffix/bogus-track/bogus-exercise-12345/subdir/numeric.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "with numeric suffix")
	})

	mux.HandleFunc("/special-char-filename#.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this is a special file")
	})

	mux.HandleFunc("/\\with-leading-backslash.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "with backslash in name")
	})
	mux.HandleFunc("/\\with\\backslashes\\in\\path.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "with backslash in path")
	})

	mux.HandleFunc("/with-leading-slash.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "this has a slash")
	})

	mux.HandleFunc("/file-3.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	})

	mux.HandleFunc("/solutions/latest", func(w http.ResponseWriter, r *http.Request) {
		team := "null"
		if teamSlug := r.FormValue("team_id"); teamSlug != "" {
			team = fmt.Sprintf(`{"name": "Bogus Team", "slug": "%s"}`, teamSlug)
		}
		payloadBody := fmt.Sprintf(payloadTemplate, requestor, team, server.URL+"/")
		fmt.Fprint(w, payloadBody)
	})
	mux.HandleFunc("/solutions/bogus-id", func(w http.ResponseWriter, r *http.Request) {
		payloadBody := fmt.Sprintf(payloadTemplate, requestor, "null", server.URL+"/")
		fmt.Fprint(w, payloadBody)
	})

	return server
}

func assertDownloadedCorrectFiles(t *testing.T, targetDir string) {
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
			desc:     "a path with a numeric suffix",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "subdir", "numeric.txt"),
			contents: "with numeric suffix",
		},
		{
			desc:     "a file that requires URL encoding",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "special-char-filename#.txt"),
			contents: "this is a special file",
		},
		{
			desc:     "a file that has a leading slash",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "with-leading-slash.txt"),
			contents: "this has a slash",
		},
		{
			desc:     "a file with a leading backslash",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "with-leading-backslash.txt"),
			contents: "with backslash in name",
		},
		{
			desc:     "a file with backslashes in path",
			path:     filepath.Join(targetDir, "bogus-track", "bogus-exercise", "with", "backslashes", "in", "path.txt"),
			contents: "with backslash in path",
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

const payloadTemplate = `
{
	"solution": {
		"id": "bogus-id",
		"user": {
			"handle": "alice",
			"is_requester": %s
		},
		"team": %s,
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
			"file-1.txt",
			"subdir/file-2.txt",
			"special-char-filename#.txt",
			"/with-leading-slash.txt",
			"\\with-leading-backslash.txt",
			"\\with\\backslashes\\in\\path.txt",
			"file-3.txt",
			"/full/path/with/numeric-suffix/bogus-track/bogus-exercise-12345/subdir/numeric.txt"
		],
		"iteration": {
			"submitted_at": "2017-08-21t10:11:12.130z"
		}
	}
}
`
