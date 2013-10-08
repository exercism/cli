package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const VERSION = "1.0.1"

var FetchEndpoints = map[string]string{
	"current": "/api/v1/user/assignments/current",
	"next":    "/api/v1/user/assignments/next",
	"demo":    "/api/v1/assignments/demo",
}

type submitResponse struct {
	Status         string
	Language       string
	Exercise       string
	SubmissionPath string `json:"submission_path"`
	Error          error
}

type submitRequest struct {
	Key  string `json:"key"`
	Code string `json:"code"`
	Path string `json:"path"`
}

func FetchAssignments(config Config, path string) (as []Assignment, err error) {
	url := fmt.Sprintf("%s%s?key=%s", config.Hostname, path, config.ApiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error fetching assignments: [%v]", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error fetching assignments. HTTP Status Code: %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		err = fmt.Errorf("Error fetching assignments: [%v]", err)
		return
	}

	var fr struct {
		Assignments []Assignment
	}

	err = json.Unmarshal(body, &fr)
	if err != nil {
		err = fmt.Errorf("Error parsing API response: [%v]", err)
		return
	}

	return fr.Assignments, err
}

func UnsubmitAssignment(config Config) (r string, err error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s?key=%s", config.Hostname, path, config.ApiKey)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", fmt.Sprintf("github.com/kytrinyx/exercism CLI v%s", VERSION))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error destroying submission: [%v]", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusNoContent {

		var ur struct {
			Error string
		}

		err = json.Unmarshal(body, &ur)
		if err != nil {
			return
		}

		err = fmt.Errorf("Status: %d, Error: %v", resp.StatusCode, ur.Error)
		return ur.Error, err
	}

	return
}

func SubmitAssignment(config Config, filePath string, code []byte) (r submitResponse, err error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s", config.Hostname, path)

	submission := submitRequest{Key: config.ApiKey, Code: string(code), Path: filePath}
	submissionJson, err := json.Marshal(submission)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(submissionJson))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", fmt.Sprintf("github.com/kytrinyx/exercism CLI v%s", VERSION))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error posting assignment: [%v]", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal(body, &r)
		if err != nil {
			return
		}
		err = fmt.Errorf("Status: %d, Error: %v", resp.StatusCode, r)
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		err = fmt.Errorf("Error parsing API response: [%v]", err)
	}

	return
}
