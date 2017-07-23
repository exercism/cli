package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	// ErrUnknownTrack represents an error returned when the track requested does not exist.
	ErrUnknownTrack = errors.New("no such track")
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
		url = fmt.Sprintf("%s/v2/exercises/%s?key=%s", c.XAPIHost, args[0], c.APIKey)
	case 2:
		url = fmt.Sprintf("%s/v2/exercises/%s/%s", c.XAPIHost, args[0], args[1])
		if c.APIKey != "" {
			url = fmt.Sprintf("%s?key=%s", url, c.APIKey)
		}
	default:
		return nil, fmt.Errorf("Usage: exercism fetch\n   or: exercism fetch TRACK_ID\n   or: exercism fetch TRACK_ID PROBLEM")
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

// FetchAll retrieves all problems for a given language track from the API
func (c *Client) FetchAll(trackID string) ([]*Problem, error) {
	list, err := c.List(trackID)
	if err != nil {
		return nil, err
	}

	problems := make([]*Problem, len(list))
	for i, prob := range list {
		p, err := c.Fetch([]string{trackID, prob})
		if err != nil {
			return nil, err
		}
		problems[i] = p[0]
	}
	return problems, nil
}

// RestoreAll fetches the latest revisions of all solutions and writes them to disk.
func (c *Client) RestoreAll() ([]*Problem, error) {
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

// Restore fetches the latest revision of specified solutions and writes them to disk.
func (c *Client) Restore(track string, exercises ...string) ([]*Problem, error) {
	result := make([]*Problem, 0, len(exercises))

	for _, exercise := range exercises {
		url := fmt.Sprintf("%s/api/v2/exercises/%s/%s?key=%s", c.XAPIHost, track, exercise, c.APIKey)
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
			log.Printf("unable to fetch problem %s/%s (HTTP: %d) - %s", track, exercise, res.StatusCode, payload.Error)
		} else {
			result = append(result, payload.Problems...)
		}
	}

	return result, nil
}

// Submissions gets a list of submitted exercises and their current state.
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

// SubmissionURL gets the url of the latest iteration on the given language track id and problem slug.
func (c *Client) SubmissionURL(trackID, slug string) (*Submission, error) {
	url := fmt.Sprintf("%s/api/v1/submissions/%s/%s?key=%s", c.APIHost, trackID, slug, c.APIKey)
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

// Submit posts an iteration to the API.
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

// List available problems for a language track.
func (c *Client) List(trackID string) ([]string, error) {
	url := fmt.Sprintf("%s/tracks/%s", c.XAPIHost, trackID)

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
		return nil, ErrUnknownTrack
	}

	var payload struct {
		Track Track
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	problems := make([]string, len(payload.Track.Problems))
	prefix := trackID + "/"

	for n, p := range payload.Track.Problems {
		problems[n] = strings.TrimPrefix(p, prefix)
	}

	return problems, nil
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

// Skip marks the exercise in the given language track as skipped.
func (c *Client) Skip(trackID, slug string) error {
	url := fmt.Sprintf("%s/api/v1/iterations/%s/%s/skip?key=%s", c.APIHost, trackID, slug, c.APIKey)

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

// Status sends a request to exercism to fetch the user's
// completion status for the given language track.
func (c *Client) Status(trackID string) (*StatusInfo, error) {
	url := fmt.Sprintf("%s/api/v1/tracks/%s/status?key=%s", c.APIHost, trackID, c.APIKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrUnknownTrack
	}

	var si StatusInfo
	if err := json.NewDecoder(res.Body).Decode(&si); err != nil {
		return nil, err
	}

	return &si, nil
}
