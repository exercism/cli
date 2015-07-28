package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	osMap = map[string]string{
		"darwin":  "mac",
		"linux":   "linux",
		"windows": "windows",
	}

	archMap = map[string]string{
		"amd64": "64bit",
		"386":   "32bit",
		"arm":   "arm",
	}
)

type asset struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

func (a *asset) Download() (*bytes.Reader, error) {
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

const installFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC

func installTgz(source *bytes.Reader, dest string) error {
	gr, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileCopy, err := os.OpenFile(dest, installFlag, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		if _, err = io.Copy(fileCopy, tr); err != nil {
			return err
		}
	}

	return nil
}

func installZip(source *bytes.Reader, dest string) error {
	zr, err := zip.NewReader(source, int64(source.Len()))
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		fileCopy, err := os.OpenFile(dest, installFlag, f.Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(fileCopy, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchLatestRelease(client http.Client) (*release, error) {
	resp, err := client.Get("https://api.github.com/repos/exercism/cli/releases/latest")
	if err != nil {
		return nil, err
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	return &rel, nil
}
