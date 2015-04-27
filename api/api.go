package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrUnknownLanguage represents an error returned when the language requested does not exist
	ErrUnknownLanguage = errors.New("the language is unknown")
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

// SubmissionInfo contains state information about a submission.
type SubmissionInfo struct {
	Slug  string `json:"slug"`
	State string `json:"state"`
}

// Fetch retrieves problems from the API.
// Most problems consist of a README, some sort of test suite, and
// any supporting files (header files, test data, boilerplate, skeleton
// files, etc).
func (c *Client) Fetch(args []string) ([]*Problem, error) {
	var url string
	switch len(args) {
	case 0:
		url = fmt.Sprintf("%s/v2/exercises?key=%s", c.XAPIHost, c.APIKey)
	case 1:
		language := args[0]
		url = fmt.Sprintf("%s/v2/exercises/%s?key=%s", c.XAPIHost, language, c.APIKey)
	case 2:
		language := args[0]
		problem := args[1]
		url = fmt.Sprintf("%s/v2/exercises/%s/%s", c.XAPIHost, language, problem)
	default:
		return nil, fmt.Errorf("Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE PROBLEM")
	}

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	payload := &PayloadProblems{}
	res, err := c.Do(req, payload)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch problems (HTTP: %d) - %s", res.StatusCode, payload.Error)
	}

	return payload.Problems, nil
}

// Restore fetches the latest revision of a solution and writes it to disk.
func (c *Client) Restore() ([]*Problem, error) {
	url := fmt.Sprintf("%s/v2/exercises/restore?key=%s", c.XAPIHost, c.APIKey)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	payload := &PayloadProblems{}
	res, err := c.Do(req, payload)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch problems (HTTP: %d) - %s", res.StatusCode, payload.Error)
	}

	return payload.Problems, nil
}

// Submissions gets a list of submitted exercises and their current acceptance state.
func (c *Client) Submissions() (map[string][]SubmissionInfo, error) {
	url := fmt.Sprintf("%s/api/v1/exercises?key=%s", c.APIHost, c.APIKey)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var payload map[string][]SubmissionInfo
	if _, err := c.Do(req, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// Submission get the latest submitted exercise for the given language and exercise
func (c *Client) Submission(language, excercise string) (*Submission, error) {
	url := fmt.Sprintf("%s/api/v1/submissions/%s/%s?key=%s", c.APIHost, language, excercise, c.APIKey)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var payload Submission
	if _, err := c.Do(req, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

// Download fetches a solution by submission key and writes it to disk.
func (c *Client) Download(submissionID string) (*Submission, error) {
	url := fmt.Sprintf("%s/api/v1/submissions/%s", c.APIHost, submissionID)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	payload := &PayloadSubmission{}
	res, err := c.Do(req, payload)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch Submission (HTTP: %d) - %s", res.StatusCode, payload.Error)
	}

	return payload.Submission, err
}

// Demo fetches the first problem in each language track.
func (c *Client) Demo() ([]*Problem, error) {
	url := fmt.Sprintf("%s/problems/demo?key=%s", c.XAPIHost, c.APIKey)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	payload := &PayloadProblems{}
	res, err := c.Do(req, payload)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch problems (HTTP: %d) - %s", res.StatusCode, payload.Error)
	}

	return payload.Problems, nil
}

// Submit posts code to the API
func (c *Client) Submit(iter *Iteration) (*Submission, error) {
	url := fmt.Sprintf("%s/api/v1/user/assignments", c.APIHost)
	payload, err := json.Marshal(iter)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	ps := &PayloadSubmission{}
	res, err := c.Do(req, ps)
	if err != nil {
		return nil, fmt.Errorf("unable to submit solution - %s", err)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unable to submit (HTTP: %d) - %s", res.StatusCode, ps.Error)
	}

	return ps.Submission, nil
}

// List available problems for a language
func (c *Client) List(language string) ([]string, error) {
	url := fmt.Sprintf("%s/tracks/%s", c.XAPIHost, language)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, ErrUnknownLanguage
	}

	var payload struct {
		Track Track
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	problems := make([]string, len(payload.Track.Problems))
	prefix := language + "/"

	for n, p := range payload.Track.Problems {
		problems[n] = strings.TrimPrefix(p, prefix)
	}

	return problems, nil
}

// Unsubmit deletes a submission.
func (c *Client) Unsubmit() error {
	url := fmt.Sprintf("%s/api/v1/user/assignments?key=%s", c.APIHost, c.APIKey)
	req, err := c.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	pe := &PayloadError{}
	if _, err := c.Do(req, pe); err != nil {
		return fmt.Errorf("failed to unsubmit - %s", pe.Error)
	}

	return nil
}

// Tracks gets the current list of active and inactive language tracks.
func (c *Client) Tracks() ([]*Track, error) {
	url := fmt.Sprintf("%s/tracks", c.XAPIHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []*Track{}, err
	}

	var payload struct {
		Tracks []*Track
	}
	if _, err := c.Do(req, &payload); err != nil {
		return []*Track{}, err
	}

	return payload.Tracks, nil
}

func (c *Client) Skip(language, slug string) error {
	url := fmt.Sprintf("%s/api/v1/iterations/%s/%s/skip?key=%s", c.APIHost, language, slug, c.APIKey)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	res, err := c.Do(req, nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	var pe PayloadError
	if err := json.NewDecoder(res.Body).Decode(&pe); err != nil {
		return err
	}

	return errors.New(pe.Error)
}
