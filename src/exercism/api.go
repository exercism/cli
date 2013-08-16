package exercism

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var FetchEndpoints = map[string]string{
	"current": "/api/v1/user/assignments/current",
	"next":    "/api/v1/user/assignments/next",
	"demo":    "/api/v1/assignments/demo",
}

type fetchResponse struct {
	Assignments []Assignment
}

type submitResponse struct {
	Status         string
	Language       string
	Exercise       string
	SubmissionPath string `json:"submission_path"`
}

func FetchAssignments(host string, path string, apiKey string) (as []Assignment, err error) {
	url := fmt.Sprintf("%s%s?key=%s", host, path, apiKey)
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

	var fr fetchResponse

	err = json.Unmarshal(body, &fr)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error parsing API response: [%s]", err.Error()))
		return
	}

	return fr.Assignments, err
}

func SubmitAssignment(host, apiKey string, a Assignment) (r *submitResponse, err error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s?key=%s", host, path, apiKey)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error posting assignment: [%s]", err.Error()))
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error posting assignment. Status: %s", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		err = errors.New(fmt.Sprintf("Error posting assignment: [%s]", err.Error()))
		return
	}

	var sr submitResponse

	err = json.Unmarshal(body, &sr)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error parsing API response: [%s]", err.Error()))
		return
	}

	return &sr, nil
}
