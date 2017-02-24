package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robphoenix/cli/config"
	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
	conf   = &config.Config{APIKey: "apikey", API: "localhost", XAPI: "xlocalhost"}
)

func TestNewRequestSetsDefaultHeaders(t *testing.T) {
	UserAgent = "Test"
	client = NewClient(conf)

	req, err := client.NewRequest("GET", client.APIHost, nil)
	assert.NoError(t, err)

	assert.Equal(t, UserAgent, req.Header.Get("User-Agent"))
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestDo(t *testing.T) {
	UserAgent = "Exercism Test v1"
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	defer server.Close()
	url := server.URL
	conf = &config.Config{APIKey: "apikey", API: url, XAPI: url}
	client = NewClient(conf)

	type test struct {
		T string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, UserAgent, r.Header.Get("User-Agent"))

		fmt.Fprint(w, `{"T":"world"}`)
	})

	req, _ := client.NewRequest("GET", client.APIHost+"/", nil)

	var body test
	_, err := client.Do(req, &body)
	assert.NoError(t, err)
	assert.Equal(t, test{T: "world"}, body)
}
