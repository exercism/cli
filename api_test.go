package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

var assignmentsJson = `
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
	fmt.Fprintf(rw, assignmentsJson)
}

func TestFetchWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fetchHandler))

	c := config.Config{
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

	c := config.Config{
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

	userAgentMatches := r.Header.Get("User-Agent") == fmt.Sprintf("github.com/exercism/cli v%s", VERSION)

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

	submitJson := `
{
	"status":"saved",
	"language":"ruby",
	"exercise":"bob",
	"submission_path":"/username/ruby/bob"
}
`
	fmt.Fprintf(rw, submitJson)
}

func TestSubmitWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))
	defer server.Close()

	var code = []byte("My source code\n")
	c := config.Config{
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

	c := config.Config{
		Hostname: server.URL,
		APIKey:   "myWrongAPIKey",
	}

	var code = []byte("My source code\n")
	_, err := SubmitAssignment(c, "ruby/bob/bob.rb", code)

	assert.Error(t, err)
}
