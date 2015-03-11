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
		APIKey   = "mykey"
		language = "go"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		languageProblemsAPI := fmt.Sprintf("/v2/exercises/%s?key=%s", language, APIKey)
		assert.Equal(t, languageProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "problems.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	_, err := client.Fetch([]string{language})
	assert.NoError(t, err)
}

func TestFetchASpecificProblem(t *testing.T) {
	var (
		APIKey   = "mykey"
		language = "go"
		problem  = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		languageProblemsAPI := fmt.Sprintf("/v2/exercises/%s/%s", language, problem)
		assert.Equal(t, languageProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "problems.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	_, err := client.Fetch([]string{language, problem})
	assert.NoError(t, err)
}

func TestGetSubmission(t *testing.T) {
	var (
		APIKey   = "mykey"
		language = "go"
		problem  = "leap"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		languageProblemsAPI := fmt.Sprintf("/api/v1/submissions/%s/%s?key=%s", language, problem, APIKey)
		assert.Equal(t, languageProblemsAPI, req.RequestURI)

		if err := respondWithFixture(w, "submission.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL, APIKey: APIKey})
	_, err := client.Submission(language, problem)
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

func TestUnknownLanguage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL})

	_, err := client.List("rubbbby")
	assert.Equal(t, err, ErrUnknownLanguage)
}
