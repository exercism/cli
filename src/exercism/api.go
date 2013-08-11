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

func Fetch(host string, apiKey string) (as []Assignment, err error) {
	path := "user/assignments/current"
	url := fmt.Sprintf("%s/%s?key=%s", host, path, apiKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusForbidden {
		err = errors.New("Unauthorized")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return
	}

	var fr FetchResponse

	err = json.Unmarshal(body, &fr)
	if err != nil {
		return
	}

	return fr.Assignments, err
}
