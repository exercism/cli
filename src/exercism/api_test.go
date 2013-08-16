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

var handler = func(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	apiKey := r.Form.Get("key")
	if r.URL.Path != "/api/v1/user/assignments/current" {
		fmt.Println("Not found")
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if apiKey != "myApiKey" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, assignmentsJson)
}

func TestFetchWithKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler))

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
	server := httptest.NewServer(http.HandlerFunc(handler))

	assignments, err := FetchAssignments(server.URL, "myWrongApiKey")

	assert.Error(t, err)
	assert.Equal(t, len(assignments), 0)

	server.Close()
}
