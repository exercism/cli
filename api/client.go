package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/exercism/cli/debug"
)

var (
	// UserAgent lets the API know where the call is being made from.
	// It's overridden from the root command so that we can set the version.
	UserAgent = "github.com/exercism/cli"

	// TimeoutInSeconds is the timeout the default HTTP client will use.
	TimeoutInSeconds = 60
	// HTTPClient is the client used to make HTTP calls in the cli package.
	HTTPClient = &http.Client{Timeout: time.Duration(TimeoutInSeconds) * time.Second}
)

// Client is an http client that is configured for Exercism.
type Client struct {
	*http.Client
	ContentType string
	Token       string
	APIBaseURL  string
}

// NewClient returns an Exercism API client.
func NewClient(token, baseURL string) (*Client, error) {
	return &Client{
		Client:     HTTPClient,
		Token:      token,
		APIBaseURL: baseURL,
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
	debug.DumpRequest(req)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	debug.DumpResponse(res)
	return res, nil
}

// TokenIsValid calls the API to determine whether the token is valid.
func (c *Client) TokenIsValid() (bool, error) {
	url := fmt.Sprintf("%s/validate_token", c.APIBaseURL)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// IsPingable calls the API /ping to determine whether the API can be reached.
func (c *Client) IsPingable() error {
	url := fmt.Sprintf("%s/ping", c.APIBaseURL)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}
	return nil
}
