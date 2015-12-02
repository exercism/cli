package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func respondWithFixture(w http.ResponseWriter, name string) error {
	f, err := os.Open("../fixtures/" + name)
	if err != nil {
		return err
	}

	io.Copy(w, f)
	f.Close()

	return nil
}
func TestFetchAllProblem(t *testing.T) {
	APIKey := "mykey"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		allProblemsAPI := fmt.Sprintf("/v2/exercises?key=%s", APIKey)
		assert.Equal(t, allProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "problems.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	problems, err := client.Fetch([]string{})
	assert.NoError(t, err)

	assert.Equal(t, len(problems), 3)
}

func TestFetchATrack(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackProblemsAPI := fmt.Sprintf("/v2/exercises/%s?key=%s", trackID, APIKey)
		assert.Equal(t, trackProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "problems.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	_, err := client.Fetch([]string{trackID})
	assert.NoError(t, err)
}

func TestFetchASpecificProblem(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
		slug    = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackProblemsAPI := fmt.Sprintf("/v2/exercises/%s/%s", trackID, slug)
		assert.Equal(t, trackProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "problems.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	_, err := client.Fetch([]string{trackID, slug})
	assert.NoError(t, err)
}

func TestSkipProblem(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
		slug    = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		skipAPI := fmt.Sprintf("/api/v1/iterations/%s/%s/skip?key=%s", trackID, slug, APIKey)
		assert.Equal(t, skipAPI, req.RequestURI)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL, APIKey: APIKey})

	err := client.Skip(trackID, slug)
	assert.NoError(t, err)
}

func TestSkipProblemErrorResponse(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
		slug    = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		skipAPI := fmt.Sprintf("/api/v1/iterations/%s/%s/skip?key=%s", trackID, slug, APIKey)
		assert.Equal(t, skipAPI, req.RequestURI)

		w.Write([]byte(`{"error":"exercise skipped"}`))
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL, APIKey: APIKey})

	err := client.Skip(trackID, slug)
	assert.Error(t, err)
}

func TestGetSubmission(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
		slug    = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackProblemsAPI := fmt.Sprintf("/api/v1/submissions/%s/%s?key=%s", trackID, slug, APIKey)
		assert.Equal(t, trackProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "submission.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL, APIKey: APIKey})
	_, err := client.IterationURL(trackID, slug)
	assert.NoError(t, err)
}

func TestSubmitAssignment(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusCreated)

		if err := respondWithFixture(w, "submit.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL})
	iter := &Iteration{} // it doesn't matter, we're testing that we can read the fixture
	sub, err := client.Submit(iter)
	assert.NoError(t, err)

	assert.Equal(t, sub.Language, "ruby")
}

func TestListTrack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// check that we correctly built the URI path
		assert.Equal(t, "/tracks/clojure", req.RequestURI)

		if err := respondWithFixture(w, "tracks.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL})

	problems, err := client.List("clojure")
	assert.NoError(t, err)

	assert.Equal(t, len(problems), 34)
	assert.Equal(t, problems[0], "bob")
}

func TestUnknownTrack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL})

	_, err := client.List("rubbbby")
	assert.Equal(t, err, ErrUnknownTrack)
}
