package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/exercism/cli/config"
)

// PayloadProblems represents a response containing problems.
type PayloadProblems struct {
	Problems []*Problem
	Error    string `json:"error"`
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

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	payload := &PayloadProblems{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing API response: [%v]", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`unable to fetch problems (HTTP: %d) - %s`, res.StatusCode, payload.Error)
	}

	return payload.Problems, nil
}

// Demo fetches the first problem in each language track.
func Demo(c *config.Config) ([]*Problem, error) {
	url := fmt.Sprintf("%s/problems/demo?key=%s", c.XAPI, c.APIKey)

	return Fetch(url)
}
