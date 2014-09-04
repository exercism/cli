package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/exercism/cli/config"
)

type PayloadProblems struct {
	Problems []*Problem
}

type PayloadError struct {
	Error string `json:"error"`
}

func Demo(c *config.Config) ([]*Problem, error) {
	url := fmt.Sprintf("%s/problems/demo?key=%s", c.ProblemsHost, c.APIKey)

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

	if res.StatusCode != http.StatusOK {
		payload := &PayloadError{}
		err := json.Unmarshal(body, payload)
		if err != nil {
			return nil, fmt.Errorf("error parsing API response: [%v]", err)
		}
		return nil, fmt.Errorf(`unable to fetch problems (%d) - %s`, res.StatusCode, payload.Error)
	}

	payload := &PayloadProblems{}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing API response: [%v]", err)
	}

	return payload.Problems, nil
}
