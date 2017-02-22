package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/robphoenix/cli/config"
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

func TestFetchAllExercises(t *testing.T) {
	APIKey := "mykey"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		allExercisesAPI := fmt.Sprintf("/v2/exercises?key=%s", APIKey)
		assert.Equal(t, allExercisesAPI, req.RequestURI)

		if err := respondWithFixture(w, "exercises.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	exercises, err := client.Fetch([]string{})
	assert.NoError(t, err)

	assert.Equal(t, 3, len(exercises))
}

func TestFetchATrack(t *testing.T) {
	var (
		APIKey  = "mykey"
		trackID = "go"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackExercisesAPI := fmt.Sprintf("/v2/exercises/%s?key=%s", trackID, APIKey)
		assert.Equal(t, trackExercisesAPI, req.RequestURI)

		if err := respondWithFixture(w, "exercises.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL, APIKey: APIKey})

	_, err := client.Fetch([]string{trackID})
	assert.NoError(t, err)
}

func TestFetchAll(t *testing.T) {
	var (
		trackID     = "go"
		fetchedList bool
		fetchCount  int
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackExercisesAPI := fmt.Sprintf("/tracks/%s", trackID)
		if req.RequestURI == trackExercisesAPI {
			if err := respondWithFixture(w, "fetch_all_tracks.json"); err != nil {
				t.Fatal(err)
			}
			fetchedList = true
		} else {
			if err := respondWithFixture(w, "exercises.json"); err != nil {
				t.Fatal(err)
			}
			fetchCount++
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL})

	all, err := client.FetchAll(trackID)

	assert.NoError(t, err)
	assert.Equal(t, len(all), 3)
	assert.Equal(t, fetchedList, true)
	assert.Equal(t, fetchCount, 3)
}

func TestFetchASpecificExercise(t *testing.T) {
	tests := []struct {
		key, url string
	}{
		{"", "/v2/exercises/go/leap"},
		{"mykey", "/v2/exercises/go/leap?key=mykey"},
	}

	for _, test := range tests {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			assert.Equal(t, test.url, req.RequestURI)

			if err := respondWithFixture(w, "exercises.json"); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		client := NewClient(&config.Config{XAPI: ts.URL, APIKey: test.key})

		_, err := client.Fetch([]string{"go", "leap"})
		assert.NoError(t, err)
	}
}

func TestSkipExercise(t *testing.T) {
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

func TestSkipExerciseErrorResponse(t *testing.T) {
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
		trackExercisesAPI := fmt.Sprintf("/api/v1/submissions/%s/%s?key=%s", trackID, slug, APIKey)
		assert.Equal(t, trackExercisesAPI, req.RequestURI)

		if err := respondWithFixture(w, "submission.json"); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL, APIKey: APIKey})
	_, err := client.SubmissionURL(trackID, slug)
	assert.NoError(t, err)
}

func TestSubmitAssignment(t *testing.T) {
	submissionComment := "hello world!"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusCreated)

		var body map[string]string
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}

		comment, ok := body["comment"]
		if ok && comment != submissionComment {
			t.Fatal("comment found and was empty")
		}

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

	// Test sending comment
	iter.Comment = submissionComment
	_, err = client.Submit(iter)
	assert.NoError(t, err)
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

	exercises, err := client.List("clojure")
	assert.NoError(t, err)

	assert.Equal(t, len(exercises), 34)
	assert.Equal(t, exercises[0], "bob")
}

func TestListUnknownTrack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	}))
	defer ts.Close()

	client := NewClient(&config.Config{XAPI: ts.URL})

	_, err := client.List("rubbbby")
	assert.Equal(t, err, ErrUnknownTrack)
}

func TestStatusUnknownTrack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	}))
	defer ts.Close()

	client := NewClient(&config.Config{API: ts.URL})

	_, err := client.Status("rubbbby")
	assert.Equal(t, err, ErrUnknownTrack)
}
