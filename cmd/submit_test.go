package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
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
	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: viper.New(),
	}

	err := runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "Welcome to Exercism", err.Error())
		assert.Regexp(t, "exercism.org/my/settings", err.Error())
	}
}

func TestSubmitWithoutWorkspace(t *testing.T) {
	v := viper.New()
	v.Set("token", "abc123")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err := runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{})
	if assert.Error(t, err) {
		assert.Regexp(t, "re-run the configure", err.Error())
	}
}

func TestSubmitNonExistentFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "submit-no-such-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://api.example.com")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err = os.WriteFile(filepath.Join(tmpDir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "file-2.txt"), []byte("This is file 2"), os.FileMode(0755))
	assert.NoError(t, err)
	files := []string{
		filepath.Join(tmpDir, "file-1.txt"),
		"no-such-file.txt",
		filepath.Join(tmpDir, "file-2.txt"),
	}
	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	if assert.Error(t, err) {
		assert.Regexp(t, "cannot be found", err.Error())
	}
}

func TestSubmitExerciseWithoutMetadataFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "no-metadata-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(file, []byte("This is a file."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://api.example.com")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file})
	if assert.Error(t, err) {
		assert.Regexp(t, "doesn't have the necessary metadata", err.Error())
	}
}

func TestGetExerciseSolutionFiles(t *testing.T) {

	tmpDir, err := os.MkdirTemp("", "dir-with-no-metadata")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	_, err = getExerciseSolutionFiles(tmpDir)
	if assert.Error(t, err) {
		assert.Regexp(t, "no files to submit", err.Error())
	}

	validTmpDir, err := os.MkdirTemp("", "dir-with-valid-metadata")
	defer os.RemoveAll(validTmpDir)
	assert.NoError(t, err)

	metadataDir := filepath.Join(validTmpDir, ".exercism")
	err = os.MkdirAll(metadataDir, os.FileMode(0755))
	assert.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(metadataDir, "config.json"),
		[]byte(`
{
	"files": {
		"solution": [
		  "expenses.go"
		]
  	}
}
`), os.FileMode(0755))
	assert.NoError(t, err)

	files, err := getExerciseSolutionFiles(validTmpDir)
	assert.NoError(t, err)
	if assert.Equal(t, len(files), 1) {
		assert.Equal(t, files[0], "expenses.go")
	}
}

func TestSubmitFilesAndDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "submit-no-such-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://api.example.com")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
		DefaultBaseURL:  "http://example.com",
	}

	err = os.WriteFile(filepath.Join(tmpDir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "file-2.txt"), []byte("This is file 2"), os.FileMode(0755))
	assert.NoError(t, err)
	files := []string{
		filepath.Join(tmpDir, "file-1.txt"),
		tmpDir,
		filepath.Join(tmpDir, "file-2.txt"),
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	if assert.Error(t, err) {
		assert.Regexp(t, "submitting a directory", err.Error())
		assert.Regexp(t, "Please change into the directory and provide the path to the file\\(s\\) you wish to submit", err.Error())
	}
}

func TestDuplicateFiles(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "duplicate-files")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1, file1})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(submittedFiles))
	assert.Equal(t, "This is file 1.", submittedFiles["file-1.txt"])
}

func TestSubmitFiles(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()
	Err = &bytes.Buffer{}

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-files")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))
	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	file2 := filepath.Join(dir, "subdir", "file-2.txt")
	err = os.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))
	assert.NoError(t, err)

	// We don't filter *.md files if you explicitly pass the file path.
	readme := filepath.Join(dir, "README.md")
	err = os.WriteFile(readme, []byte("This is the readme."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	files := []string{
		file1, file2, readme,
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)

	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(submittedFiles))
		assert.Equal(t, "This is file 1.", submittedFiles["file-1.txt"])
		assert.Equal(t, "This is file 2.", submittedFiles["subdir/file-2.txt"])
		assert.Equal(t, "This is the readme.", submittedFiles["README.md"])
		assert.Regexp(t, "submitted successfully", Err)
	}
}

func TestLegacyMetadataMigration(t *testing.T) {
	co := newCapturedOutput()
	co.newErr = &bytes.Buffer{}
	co.override()
	defer co.reset()

	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "legacy-metadata-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	metadata := &workspace.ExerciseMetadata{
		ID:           "bogus-solution-uuid",
		Track:        "bogus-track",
		ExerciseSlug: "bogus-exercise",
		URL:          "http://example.com/bogus-url",
		IsRequester:  true,
	}
	b, err := json.Marshal(metadata)
	assert.NoError(t, err)
	exercise := workspace.NewExerciseFromDir(dir)
	err = os.WriteFile(exercise.LegacyMetadataFilepath(), b, os.FileMode(0600))
	assert.NoError(t, err)

	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(file, []byte("This is a file."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)
	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	ok, _ := exercise.HasLegacyMetadata()
	assert.True(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.False(t, ok)

	flags := pflag.NewFlagSet("fake", pflag.PanicOnError)
	flags.Bool("verbose", true, "")

	err = runSubmit(cfg, flags, []string{file})
	assert.NoError(t, err)
	assert.Equal(t, "This is a file.", submittedFiles["file.txt"])

	ok, _ = exercise.HasLegacyMetadata()
	assert.False(t, ok)
	ok, _ = exercise.HasMetadata()
	assert.True(t, ok)
	assert.Regexp(t, "Migrated metadata", Err)
}

func TestSubmitWithEmptyFile(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "empty-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte(""), os.FileMode(0755))
	file2 := filepath.Join(dir, "file-2.txt")
	err = os.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1, file2})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(submittedFiles))
	assert.Equal(t, "This is file 2.", submittedFiles["file-2.txt"])
}

func TestSubmitWithEnormousFile(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "enormous-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(file, make([]byte, 65535), os.FileMode(0755))
	if err != nil {
		t.Fatal(err)
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file})

	if assert.Error(t, err) {
		assert.Regexp(t, "Please reduce the size of the file and try again.", err.Error())
	}
}

func TestSubmitFilesForTeamExercise(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-files")
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "teams", "bogus-team", "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))
	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	file2 := filepath.Join(dir, "subdir", "file-2.txt")
	err = os.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	files := []string{
		file1, file2,
	}
	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(submittedFiles))

	assert.Equal(t, "This is file 1.", submittedFiles["file-1.txt"])
	assert.Equal(t, "This is file 2.", submittedFiles["subdir/file-2.txt"])
}

func TestSubmitOnlyEmptyFile(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	tmpDir, err := os.MkdirTemp("", "just-an-empty-file")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://api.example.com")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(file, []byte(""), os.FileMode(0755))

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file})
	if assert.Error(t, err) {
		assert.Regexp(t, "No files found", err.Error())
	}
}

func TestSubmitFilesFromDifferentSolutions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dir-1-submit")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir1 := filepath.Join(tmpDir, "bogus-track", "bogus-exercise-1")
	os.MkdirAll(dir1, os.FileMode(0755))
	writeFakeMetadata(t, dir1, "bogus-track", "bogus-exercise-1")

	dir2 := filepath.Join(tmpDir, "bogus-track", "bogus-exercise-2")
	os.MkdirAll(dir2, os.FileMode(0755))
	writeFakeMetadata(t, dir2, "bogus-track", "bogus-exercise-2")

	file1 := filepath.Join(dir1, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	file2 := filepath.Join(dir2, "file-2.txt")
	err = os.WriteFile(file2, []byte("This is file 2."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", "http://api.example.com")

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1, file2})
	if assert.Error(t, err) {
		assert.Regexp(t, "different solutions", err.Error())
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
			body, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}
			// Following RFC 7578, Go 1.17+ strips the directory information in fileHeader.Filename.
			// Validating the submitted files directory tree is important so Content-Disposition is used for
			// obtaining the unmodified filename.
			v := fileHeader.Header.Get("Content-Disposition")
			_, dispositionParams, err := mime.ParseMediaType(v)
			if err != nil {
				t.Fatalf("failed to obtain submitted filename from multipart header: %s", err.Error())
			}
			filename := dispositionParams["filename"]
			submittedFiles[filename] = string(body)
		}

		fmt.Fprint(w, "{}")
	})
	return httptest.NewServer(handler)
}

func TestSubmitRelativePath(t *testing.T) {
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	// The fake endpoint will populate this when it receives the call from the command.
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "relative-path")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(dir, os.FileMode(0755))

	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		UserViperConfig: v,
	}

	err = os.WriteFile(filepath.Join(dir, "file.txt"), []byte("This is a file."), os.FileMode(0755))

	err = os.Chdir(dir)
	assert.NoError(t, err)

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{"file.txt"})
	assert.NoError(t, err)

	assert.Equal(t, 1, len(submittedFiles))
	assert.Equal(t, "This is a file.", submittedFiles["file.txt"])
}

func TestSubmitServerErr(t *testing.T) {
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

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))
	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	err = os.WriteFile(filepath.Join(dir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	files := []string{
		filepath.Join(dir, "file-1.txt"),
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)

	assert.Regexp(t, "test error", err.Error())
}

func TestHandleErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-nonsuccess")
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

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))
	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	err = os.WriteFile(filepath.Join(dir, "file-1.txt"), []byte("This is file 1"), os.FileMode(0755))
	assert.NoError(t, err)

	files := []string{
		filepath.Join(dir, "file-1.txt"),
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), files)
	assert.Error(t, err)
}

func TestSubmissionNotConnectedToRequesterAccount(t *testing.T) {
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-files")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))

	metadata := &workspace.ExerciseMetadata{
		ID:           "bogus-solution-uuid",
		Track:        "bogus-track",
		ExerciseSlug: "bogus-exercise",
		URL:          "http://example.com/bogus-url",
		IsRequester:  false,
	}
	err = metadata.Write(dir)
	assert.NoError(t, err)

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1})
	if assert.Error(t, err) {
		assert.Regexp(t, "not connected to your account", err.Error())
	}
}

func TestExerciseDirnameMatchesMetadataSlug(t *testing.T) {
	submittedFiles := map[string]string{}
	ts := fakeSubmitServer(t, submittedFiles)
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "submit-files")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	dir := filepath.Join(tmpDir, "bogus-track", "bogus-exercise-doesnt-match-metadata-slug")
	os.MkdirAll(filepath.Join(dir, "subdir"), os.FileMode(0755))
	writeFakeMetadata(t, dir, "bogus-track", "bogus-exercise")

	file1 := filepath.Join(dir, "file-1.txt")
	err = os.WriteFile(file1, []byte("This is file 1."), os.FileMode(0755))
	assert.NoError(t, err)

	v := viper.New()
	v.Set("token", "abc123")
	v.Set("workspace", tmpDir)
	v.Set("apibaseurl", ts.URL)

	cfg := config.Config{
		Persister:       config.InMemoryPersister{},
		Dir:             tmpDir,
		UserViperConfig: v,
	}

	err = runSubmit(cfg, pflag.NewFlagSet("fake", pflag.PanicOnError), []string{file1})
	if assert.Error(t, err) {
		assert.Regexp(t, "directory does not match exercise slug", err.Error())
	}
}

func writeFakeMetadata(t *testing.T, dir, trackID, exerciseSlug string) {
	metadata := &workspace.ExerciseMetadata{
		ID:           "bogus-solution-uuid",
		Track:        trackID,
		ExerciseSlug: exerciseSlug,
		URL:          "http://example.com/bogus-url",
		IsRequester:  true,
	}
	err := metadata.Write(dir)
	assert.NoError(t, err)
}
