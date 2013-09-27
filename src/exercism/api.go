package exercism

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const VERSION = "1.0.0"

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

type submitError struct {
	Error string
}

type submitRequest struct {
	Key  string `json:"key"`
	Code string `json:"code"`
	Path string `json:"path"`
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
		err = errors.New(fmt.Sprintf("Error fetching assignments. HTTP Status Code: %d", resp.StatusCode))
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

func SubmitAssignment(host, apiKey, filePath string, code []byte) (r *submitResponse, err error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s", host, path)

	submission := submitRequest{Key: apiKey, Code: string(code), Path: filePath}
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
		err = errors.New(fmt.Sprintf("Error posting assignment: [%s]", err.Error()))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		postError := submitError{}
		_ = json.Unmarshal(body, &postError)
		err = errors.New(fmt.Sprintf("Status: %d, Error: %s", resp.StatusCode, postError.Error))
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error parsing API response: [%s]", err.Error()))
	}

	return
}
