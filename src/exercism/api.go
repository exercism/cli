package exercism

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FetchResponse struct {
	Assignments []Assignment
}

func FetchAssignments(host string, apiKey string) (as []Assignment, err error) {
	path := "api/v1/user/assignments/current"
	url := fmt.Sprintf("%s/%s?key=%s", host, path, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error fetching assignments: [%s]", err.Error()))
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Error fecthing assignments. HTTP Status Code: %d", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		err = errors.New(fmt.Sprintf("Error fetching assignments: [%s]", err.Error()))
		return
	}

	var fr FetchResponse

	err = json.Unmarshal(body, &fr)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error parsing API response: [%s]", err.Error()))
		return
	}

	return fr.Assignments, err
}
