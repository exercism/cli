package exercism

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var assignmentsJson = `
{
    "assignments": [
        {
            "track": "ruby",
            "slug": "bob",
            "readme": "Readme text",
            "test_file": "bob_test.rb",
            "tests": "Tests Text"
        }
    ]
}
`

var submitJson = `
{
	"status":"saved",
	"language":"ruby",
	"exercise":"bob",
	"submission_path":"/username/ruby/bob"
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
	if apiKey != "myApiKey" {
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, `{"error": "Unable to identify user"}`)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, assignmentsJson)
}

var submitHandler = func(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	apiKey := r.Form.Get("key")
	if r.URL.Path != "/api/v1/user/assignments" || r.Method != "POST" {
		fmt.Println("Not found")
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if apiKey != "myApiKey" {
		rw.WriteHeader(http.StatusForbidden)
		rw.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(rw, `{"error": "Unable to identify user"}`)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(rw, submitJson)
}

func TestFetchWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fetchHandler))

	assignments, err := FetchAssignments(server.URL, "myApiKey")
	assert.NoError(t, err)

	assert.Equal(t, len(assignments), 1)

	assert.Equal(t, assignments[0].Track, "ruby")
	assert.Equal(t, assignments[0].Slug, "bob")
	assert.Equal(t, assignments[0].Readme, "Readme text")
	assert.Equal(t, assignments[0].TestFile, "bob_test.rb")
	assert.Equal(t, assignments[0].Tests, "Tests Text")

	server.Close()
}

func TestFetchWithIncorrectKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fetchHandler))

	assignments, err := FetchAssignments(server.URL, "myWrongApiKey")

	assert.Error(t, err)
	assert.Equal(t, len(assignments), 0)

	server.Close()
}

func TestSubmitWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(submitHandler))

	assignment := Assignment{}
	response, err := SubmitAssignment(server.URL, "myApiKey", assignment)
	assert.NoError(t, err)

	assert.Equal(t, response.Status, "saved")
	assert.Equal(t, response.Language, "ruby")
	assert.Equal(t, response.Exercise, "bob")
	assert.Equal(t, response.SubmissionPath, "/username/ruby/bob")

	server.Close()
}
