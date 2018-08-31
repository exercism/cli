package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// UserAgent lets the API know where the call is being made from.
	// It's overridden from the root command so that we can set the version.
	UserAgent = "github.com/exercism/cli"

	// HTTPClient is the client used to make HTTP calls in the cli package.
	HTTPClient = &http.Client{Timeout: 60 * time.Second}
)

// Client is an http client that is configured for Exercism.
type Client struct {
	APIBaseURL  string
	ContentType string
	Token       string
	*http.Client
}

// NewClient returns an Exercism API client.
func NewClient(token, baseURL string) (*Client, error) {
	return &Client{
		APIBaseURL: baseURL,
		Token:      token,
		Client:     HTTPClient,
	}, nil
}

// NewRequest returns an http.Request with information for the Exercism API.
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	if c.Client == nil {
		c.Client = HTTPClient
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	if c.ContentType == "" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", c.ContentType)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	return req, nil
}

// Do performs an http.Request and optionally parses the response body into the given interface.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	DumpRequest(req)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	DumpResponse(res)
	return res, nil
}

// ValidateToken calls the API to determine whether the token is valid.
func (c *Client) ValidateToken() (bool, error) {
	url := fmt.Sprintf("%s/validate_token", c.APIBaseURL)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}

// Ping calls the API /ping to determine whether the API can be reached.
func (c *Client) Ping() error {
	url := fmt.Sprintf("%s/ping", c.APIBaseURL)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}
	return nil
}
