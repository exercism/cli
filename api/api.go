package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/exercism/cli/config"
)

var (
	// UserAgent lets the API know where the call is being made from.
	// It's set from main() so that we have access to the version.
	UserAgent string
)

// PayloadError represents an error message from the API.
type PayloadError struct {
	Error string `json:"error"`
}

// PayloadProblems represents a response containing problems.
type PayloadProblems struct {
	Problems []*Problem
	PayloadError
}

// PayloadSubmission represents metadata about a successful submission.
type PayloadSubmission struct {
	*Submission
	PayloadError
}

// Fetch retrieves problems from the API.
// In most cases these problems consist of a test suite and a README
// from the x-api, but it is also used when restoring earlier iterations.
func Fetch(url string) ([]*Problem, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	payload := &PayloadProblems{}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(payload); err != nil {
		return nil, fmt.Errorf("error parsing API response - %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`unable to fetch problems (HTTP: %d) - %s`, res.StatusCode, payload.Error)
	}

	return payload.Problems, nil
}

func Download(url string) (*Submission, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	payload := &PayloadSubmission{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing API response - %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`unable to fetch Submission (HTTP: %d) - %s`, res.StatusCode, payload.Error)
	}

	return payload.Submission, err
}

// Demo fetches the first problem in each language track.
func Demo(c *config.Config) ([]*Problem, error) {
	url := fmt.Sprintf("%s/problems/demo?key=%s", c.XAPI, c.APIKey)

	return Fetch(url)
}

// Submit posts code to the API
func Submit(url string, iter *Iteration) (*Submission, error) {
	payload, err := json.Marshal(iter)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to submit solution - %s", err)
	}
	defer res.Body.Close()

	ps := &PayloadSubmission{}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(ps); err != nil {
		return nil, fmt.Errorf("error parsing API response - %s", err)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(`unable to submit (HTTP: %d) - %s`, res.StatusCode, ps.Error)
	}

	return ps.Submission, nil
}

// Unsubmit deletes a submission.
func Unsubmit(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	pe := &PayloadError{}
	if err := json.NewDecoder(res.Body).Decode(pe); err != nil {
		return fmt.Errorf("failed to unsubmit - %s", err)
	}
	return fmt.Errorf("failed to unsubmit - %s", pe.Error)
}

// Tracks gets the current list of active and inactive language tracks.
func Tracks(url string) ([]*Track, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []*Track{}, err
	}
	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []*Track{}, err
	}
	defer res.Body.Close()

	var payload struct {
		Tracks []*Track
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&payload)
	if err != nil {
		return []*Track{}, err
	}
	return payload.Tracks, nil
}
