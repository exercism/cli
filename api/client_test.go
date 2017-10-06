package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exercism/cli/config"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestSetsDefaultHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `ok`)
	}))
	defer ts.Close()

	UserAgent = "BogusAgent"

	tests := []struct {
		client      *Client
		auth        string
		contentType string
	}{
		{
			// Use defaults.
			client:      &Client{},
			auth:        "",
			contentType: "application/json",
		},
		{
			// Override defaults.
			client: &Client{
				UserConfig:  &config.UserConfig{Token: "abc123"},
				ContentType: "bogus",
			},
			auth:        "Bearer abc123",
			contentType: "bogus",
		},
	}

	for _, test := range tests {
		req, err := test.client.NewRequest("GET", ts.URL, nil)
		assert.NoError(t, err)
		assert.Equal(t, "BogusAgent", req.Header.Get("User-Agent"))
		assert.Equal(t, test.contentType, req.Header.Get("Content-Type"))
		assert.Equal(t, test.auth, req.Header.Get("Authorization"))
	}
}

func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		fmt.Fprint(w, `{"hello": "world"}`)
	}))
	defer ts.Close()

	type payload struct {
		Hello string `json:"hello"`
	}

	client := &Client{}

	req, err := client.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)

	var body payload
	_, err = client.Do(req, &body)
	assert.NoError(t, err)
	assert.Equal(t, "world", body.Hello)
}
