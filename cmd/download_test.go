package cmd

import (
	"encoding/json"
	"fmt"
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
		assert.Regexp(t, "exercism.org/my/settings", err.Error())
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

func TestSolutionFile(t *testing.T) {
	testCases := []struct {
		name, file, expectedPath, expectedURL string
	}{
		{
			name:         "filename with special character",
			file:         "special-char-filename#.txt",
			expectedPath: "special-char-filename#.txt",
			expectedURL:  "http://www.example.com/special-char-filename%23.txt",
		},
		{
			name:         "filename with leading slash",
			file:         "/with-leading-slash.txt",
			expectedPath: fmt.Sprintf("%cwith-leading-slash.txt", os.PathSeparator),
			expectedURL:  "http://www.example.com//with-leading-slash.txt",
		},
		{
			name:         "filename with leading backslash",
			file:         "\\with-leading-backslash.txt",
			expectedPath: fmt.Sprintf("%cwith-leading-backslash.txt", os.PathSeparator),
			expectedURL:  "http://www.example.com/%5Cwith-leading-backslash.txt",
		},
		{
			name:         "filename with backslashes in path",
			file:         "\\backslashes\\in-path.txt",
			expectedPath: fmt.Sprintf("%[1]cbackslashes%[1]cin-path.txt", os.PathSeparator),
			expectedURL:  "http://www.example.com/%5Cbackslashes%5Cin-path.txt",
		},
		{
			name:         "path with a numeric suffix",
			file:         "/bogus-exercise-12345/numeric.txt",
			expectedPath: fmt.Sprintf("%[1]cbogus-exercise-12345%[1]cnumeric.txt", os.PathSeparator),
			expectedURL:  "http://www.example.com//bogus-exercise-12345/numeric.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sf := solutionFile{
				path:    tc.file,
				baseURL: "http://www.example.com/",
			}

			if sf.relativePath() != tc.expectedPath {
				t.Fatalf("Expected path '%s', got '%s'", tc.expectedPath, sf.relativePath())
			}

			url, err := sf.url()
			if err != nil {
				t.Fatal(err)
			}

			if url != tc.expectedURL {
				t.Fatalf("Expected URL '%s', got '%s'", tc.expectedURL, url)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

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
		tmpDir, err := os.MkdirTemp("", "download-cmd")
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
		b, err := os.ReadFile(workspace.NewExerciseFromDir(dir).MetadataFilepath())
		assert.NoError(t, err)
		var metadata workspace.ExerciseMetadata
		err = json.Unmarshal(b, &metadata)
		assert.NoError(t, err)

		assert.Equal(t, "bogus-track", metadata.Track)
		assert.Equal(t, "bogus-exercise", metadata.ExerciseSlug)
		assert.Equal(t, tc.requester, metadata.IsRequester)
	}
}

func TestDownloadToExistingDirectory(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	testCases := []struct {
		exerciseDir string
		flags       map[string]string
	}{
		{
			exerciseDir: filepath.Join("bogus-track", "bogus-exercise"),
			flags:       map[string]string{"exercise": "bogus-exercise", "track": "bogus-track"},
		},
		{
			exerciseDir: filepath.Join("teams", "bogus-team", "bogus-track", "bogus-exercise"),
			flags:       map[string]string{"exercise": "bogus-exercise", "track": "bogus-track", "team": "bogus-team"},
		},
	}

	for _, tc := range testCases {
		tmpDir, err := os.MkdirTemp("", "download-cmd")
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		err = os.MkdirAll(filepath.Join(tmpDir, tc.exerciseDir), os.FileMode(0755))
		assert.NoError(t, err)

		ts := fakeDownloadServer("true", "")
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

		if assert.Error(t, err) {
			assert.Regexp(t, "directory '.+' already exists", err.Error())
		}
	}
}

func TestDownloadToExistingDirectoryWithForce(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	testCases := []struct {
		exerciseDir string
		flags       map[string]string
	}{
		{
			exerciseDir: filepath.Join("bogus-track", "bogus-exercise"),
			flags:       map[string]string{"exercise": "bogus-exercise", "track": "bogus-track"},
		},
		{
			exerciseDir: filepath.Join("teams", "bogus-team", "bogus-track", "bogus-exercise"),
			flags:       map[string]string{"exercise": "bogus-exercise", "track": "bogus-track", "team": "bogus-team"},
		},
	}

	for _, tc := range testCases {
		tmpDir, err := os.MkdirTemp("", "download-cmd")
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		err = os.MkdirAll(filepath.Join(tmpDir, tc.exerciseDir), os.FileMode(0755))
		assert.NoError(t, err)

		ts := fakeDownloadServer("true", "")
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
		flags.Set("force", "true")

		err = runDownload(cfg, flags, []string{})
		assert.NoError(t, err)
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
	}

	for _, file := range expectedFiles {
		t.Run(file.desc, func(t *testing.T) {
			b, err := os.ReadFile(file.path)
			assert.NoError(t, err)
			assert.Equal(t, file.contents, string(b))
		})
	}

	path := filepath.Join(targetDir, "bogus-track", "bogus-exercise", "file-3.txt")
	_, err := os.Lstat(path)
	assert.True(t, os.IsNotExist(err), "It should not write the file if empty.")
}

func TestDownloadError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": {"type": "error", "message": "test error"}}`)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-err-tmp-dir")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	setupDownloadFlags(flags)
	flags.Set("uuid", "value")

	err = runDownload(cfg, flags, []string{})

	assert.Equal(t, "test error", err.Error())

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
			"file-3.txt"
		],
		"iteration": {
			"submitted_at": "2017-08-21t10:11:12.130z"
		}
	}
}
`
