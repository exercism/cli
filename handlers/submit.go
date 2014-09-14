package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"io/ioutil"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

const (
	msgPleaseAuthenticate = "You must be authenticated. Run `exercism configure --key=YOUR_API_KEY`."
)

var testExtensions = map[string]string{
	"ruby":    "_test.rb",
	"js":      ".spec.js",
	"elixir":  "_test.exs",
	"clojure": "_test.clj",
	"python":  "_test.py",
	"go":      "_test.go",
	"haskell": "_test.hs",
	"cpp":     "_test.cpp",
}

func Submit(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		fmt.Println("Please enter a file name")
		return
	}

	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if !c.IsAuthenticated() {
		fmt.Println(msgPleaseAuthenticate)
		return
	}

	filename := ctx.Args()[0]

	// Make filename relative to config.Dir.
	absPath, err := absolutePath(filename)
	if err != nil {
		fmt.Printf("Couldn't find %v: %v\n", filename, err)
		return
	}
	exDir := c.Dir + string(filepath.Separator)
	if !strings.HasPrefix(absPath, exDir) {
		fmt.Printf("%v is not under your exercism project path (%v)\n", absPath, exDir)
		return
	}
	filename = absPath[len(exDir):]

	if IsTest(filename) {
		fmt.Println("It looks like this is a test, please submit a solution.")
		return
	}

	code, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading %v: %v\n", absPath, err)
		return
	}

	response, err := SubmitAssignment(c, filename, code)
	if err != nil {
		fmt.Printf("There was an issue with your submission: %v\n", err)
		return
	}

	fmt.Printf("For feedback on your submission visit %s%s%s\n",
		c.Hostname, "/submissions/", response.ID)
}

type submitResponse struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Language       string `json:"language"`
	Exercise       string `json:"exercise"`
	SubmissionPath string `json:"submission_path"`
	Error          string `json:"error"`
}

type submitRequest struct {
	Key  string `json:"key"`
	Code string `json:"code"`
	Path string `json:"path"`
}

func SubmitAssignment(c *config.Config, filePath string, code []byte) (*submitResponse, error) {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s", c.Hostname, path)

	submission := submitRequest{Key: c.APIKey, Code: string(code), Path: filePath}
	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(submissionJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error posting assignment: [%v]", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var r submitResponse
	if resp.StatusCode != http.StatusCreated {
		err = json.Unmarshal(body, &r)
		if err != nil {
			return nil, err
		}
		err = fmt.Errorf("Status: %d, Error: %v", resp.StatusCode, r)
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("Error parsing API response: [%v]", err)
	}
	return &r, nil
}

func absolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

func IsTest(filename string) bool {
	for _, ext := range testExtensions {
		if strings.LastIndex(filename, ext) > 0 {
			return true
		}
	}
	return false
}
