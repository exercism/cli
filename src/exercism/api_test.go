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
	if r.URL.Path != "/user/assignments/current" {
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

	assignments, err := Fetch(server.URL, "myApiKey")
	assert.NoError(t, err)

	assert.Equal(t, len(assignments), 1)

	server.Close()
}

func TestFetchWithIncorrectKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler))

	assignments, err := Fetch(server.URL, "myWrongApiKey")

	assert.Error(t, err, "Unauthorized")
	assert.Equal(t, len(assignments), 0)

	server.Close()
}
