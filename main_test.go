package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

var assignmentsJSON = `
{
    "assignments": [
        {
            "track": "ruby",
            "slug": "bob",
						"files": {
							"README.md": "Readme text",
							"bob_test.rb": "Tests text"
						}
        }
    ]
}
`

var fetchHandler = func(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	apiKey := r.Form.Get("key")
	if r.URL.Path != "/api/v1/user/assignments/current" {
		fmt.Println("Not found")
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if apiKey != "myAPIKey" {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(rw, `{"error": "Unable to identify user"}`)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, assignmentsJSON)
}

func TestFetchWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fetchHandler))

	c := &config.Config{
		Hostname: server.URL,
		APIKey:   "myAPIKey",
	}

	assignments, err := FetchAssignments(c, "/api/v1/user/assignments/current")
	assert.NoError(t, err)

	assert.Equal(t, len(assignments), 1)

	assert.Equal(t, assignments[0].Track, "ruby")
	assert.Equal(t, assignments[0].Slug, "bob")
	assert.Equal(t, assignments[0].Files, map[string]string{
		"README.md":   "Readme text",
		"bob_test.rb": "Tests text",
	},
	)

	server.Close()
}

func TestFetchWithIncorrectKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fetchHandler))

	c := &config.Config{
		Hostname: server.URL,
		APIKey:   "myWrongAPIKey",
	}

	assignments, err := FetchAssignments(c, "/api/v1/user/assignments/current")

	assert.Error(t, err)
	assert.Equal(t, len(assignments), 0)
	assert.Contains(t, fmt.Sprintf("%s", err), "Unable to identify user")

	server.Close()
}

var submitHandler = func(rw http.ResponseWriter, r *http.Request) {
	pathMatches := r.URL.Path == "/api/v1/user/assignments"
	methodMatches := r.Method == "POST"
	if !(pathMatches && methodMatches) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	userAgentMatches := strings.HasPrefix(r.Header.Get("User-Agent"), fmt.Sprintf("github.com/exercism/cli v%s", Version))

	if !userAgentMatches {
		fmt.Printf("User agent mismatch: %s\n", r.Header.Get("User-Agent"))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Reading body error: %s\n", err)
		return
	}

	type Submission struct {
		Key  string
		Code string
		Path string
	}

	submission := Submission{}

	err = json.Unmarshal(body, &submission)
	if err != nil {
		fmt.Printf("Unmarshalling error: %v, Body: %s\n", err, body)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if submission.Key != "myAPIKey" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(rw, `{"error": "Unable to identify user"}`)
		return
	}

	code := submission.Code
	filePath := submission.Path

	codeMatches := string(code) == "My source code\n"
	filePathMatches := filePath == "ruby/bob/bob.rb"

	if !filePathMatches {
		fmt.Printf("FilePathMismatch: File Path: %s\n", filePath)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !codeMatches {
		fmt.Printf("Code Mismatch: Code: %v\n", code)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	submitJSON := `
{
	"status":"saved",
	"language":"ruby",
	"exercise":"bob",
	"submission_path":"/username/ruby/bob"
}
`
	fmt.Fprintf(rw, submitJSON)
}

func TestSubmitWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))
	defer server.Close()

	var code = []byte("My source code\n")
	c := &config.Config{
		Hostname: server.URL,
		APIKey:   "myAPIKey",
	}
	response, err := SubmitAssignment(c, "ruby/bob/bob.rb", code)
	assert.NoError(t, err)

	assert.Equal(t, response.Status, "saved")
	assert.Equal(t, response.Language, "ruby")
	assert.Equal(t, response.Exercise, "bob")
	assert.Equal(t, response.SubmissionPath, "/username/ruby/bob")
}

func TestSubmitWithIncorrectKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))
	defer server.Close()

	c := &config.Config{
		Hostname: server.URL,
		APIKey:   "myWrongAPIKey",
	}

	var code = []byte("My source code\n")
	_, err := SubmitAssignment(c, "ruby/bob/bob.rb", code)

	assert.Error(t, err)
}

func TestSavingAssignment(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	prepareFixture(t, fmt.Sprintf("%s/ruby/bob/stub.rb", tmpDir), "Existing stub")

	assignment := Assignment{
		Track: "ruby",
		Slug:  "bob",
		Files: map[string]string{
			"bob_test.rb":     "Tests text",
			"README.md":       "Readme text",
			"path/to/file.rb": "File text",
			"stub.rb":         "New version of stub",
		},
	}

	err = SaveAssignment(tmpDir, assignment)
	assert.NoError(t, err)

	readme, err := ioutil.ReadFile(tmpDir + "/ruby/bob/README.md")
	assert.NoError(t, err)
	assert.Equal(t, string(readme), "Readme text")

	tests, err := ioutil.ReadFile(tmpDir + "/ruby/bob/bob_test.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(tests), "Tests text")

	fileInDir, err := ioutil.ReadFile(tmpDir + "/ruby/bob/path/to/file.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(fileInDir), "File text")

	stubFile, err := ioutil.ReadFile(tmpDir + "/ruby/bob/stub.rb")
	assert.NoError(t, err)
	assert.Equal(t, string(stubFile), "Existing stub")
}

func prepareFixture(t *testing.T, fixture, s string) {
	err := os.MkdirAll(filepath.Dir(fixture), 0755)
	assert.NoError(t, err)

	err = ioutil.WriteFile(fixture, []byte(s), 0644)
	assert.NoError(t, err)

	// ensure fixture is set up correctly
	fixtureContents, err := ioutil.ReadFile(fixture)
	assert.NoError(t, err)
	assert.Equal(t, string(fixtureContents), s)
}

func TestFetchCurrentEndpoint(t *testing.T) {
	expected := "/api/v1/user/assignments/current"
	actual := FetchEndpoint([]string{})
	assert.Equal(t, expected, actual)
}

func TestFetchExerciseEndpoint(t *testing.T) {
	expected := "/api/v1/assignments/language/slug"
	actual := FetchEndpoint([]string{"language", "slug"})
	assert.Equal(t, expected, actual)
}

func TestFetchExerciseEndpointByLanguage(t *testing.T) {
	expected := "/api/v1/assignments/language"
	actual := FetchEndpoint([]string{"language"})
	assert.Equal(t, expected, actual)
}
