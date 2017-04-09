package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type asset struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

func (a *asset) download() (*bytes.Reader, error) {
	downloadURL := fmt.Sprintf("https://api.github.com/repos/exercism/cli/releases/assets/%d", a.ID)
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, err
	}
	// https://developer.github.com/v3/repos/releases/#get-a-single-release-asset
	req.Header.Set("Accept", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bs), nil
}

type release struct {
	Location string  `json:"html_url"`
	TagName  string  `json:"tag_name"`
	Assets   []asset `json:"assets"`
}

func (r *release) Version() string {
	return strings.TrimPrefix(r.TagName, "v")
}
