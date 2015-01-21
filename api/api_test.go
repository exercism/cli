package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestListTrack(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that we correctly built the URI path
		assert.Equal(t, "/tracks/clojure", r.RequestURI)

		f, err := os.Open("../fixtures/tracks.json")
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(w, f)
		f.Close()
	}))
	defer ts.Close()

	conf := &config.Config{XAPI: ts.URL}
	client := NewClient(conf)

	problems, err := client.List("clojure")
	assert.NoError(t, err)

	assert.Equal(t, len(problems), 34)
	assert.Equal(t, problems[0], "bob")
}

func TestUnknownLanguage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer ts.Close()

	conf := &config.Config{XAPI: ts.URL}
	client := NewClient(conf)

	_, err := client.List("rubbbby")
	assert.Equal(t, err, ErrUnknownLanguage)
}
