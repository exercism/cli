package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/handlers"
)

const (
	// Version is the current release of the command-line app.
	// We try to follow Semantic Versioning (http://semver.org),
	// but with the http://exercism.io app being a prototype, a
	// lot of things get out of hand.
	Version = "1.6.2"
	// UserAgent is sent along as a header to HTTP requests that the
	// CLI makes. This helps with debugging.
	UserAgent = "github.com/exercism/cli v" + Version
)

var FetchEndpoints = map[string]string{
	"current":  "/api/v1/user/assignments/current",
	"next":     "/api/v1/user/assignments/next",
	"restore":  "/api/v1/user/assignments/restore",
	"demo":     "/api/v1/assignments/demo",
	"exercise": "/api/v1/assignments",
}

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

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "configure",
			Usage: "Write config values to a JSON file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "path to exercises directory",
				},
				cli.StringFlag{
					Name:  "host, u",
					Usage: "exercism api host",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "exercism.io API key (see http://exercism.io/account)",
				},
			},
			Action: handlers.Configure,
		},
		{
			Name:      "current",
			ShortName: "c",
			Usage:     "Show the current assignments",
			Action: func(ctx *cli.Context) {
				var language string
				argc := len(ctx.Args())
				if argc != 0 && argc != 1 {
					fmt.Println("Usage: exercism current\n   or: exercism current LANGUAGE")
					return
				}

				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
					return
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				currentAssignments, err := FetchAssignments(c, FetchEndpoints["current"])
				if err != nil {
					fmt.Println(err)
					return
				}

				if argc == 1 {
					language = ctx.Args()[0]
					fmt.Println("Current Assignments for", strings.Title(language))
				} else {
					fmt.Println("Current Assignments")
				}

				for _, a := range currentAssignments {
					if argc == 1 {
						if strings.ToLower(language) == strings.ToLower(a.Track) {
							fmt.Printf("%v: %v\n", strings.Title(a.Track), a.Slug)
						}
					} else {
						fmt.Printf("%v: %v\n", strings.Title(a.Track), a.Slug)
					}
				}
			},
		},
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     "Fetch first assignment for each language from exercism.io",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
					return
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					c = config.Demo()
				}
				assignments, err := FetchAssignments(c, FetchEndpoints["demo"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(c.Dir, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				msg := "\nThe demo exercises have been written to %s, in subdirectories by language.\n\nTo try an exercise, change directory to a language/exercise, read the README and run the tests.\n\n"
				fmt.Printf(msg, c.Dir)
			},
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     "Fetch assignments from exercism.io",
			Action: func(ctx *cli.Context) {
				argCount := len(ctx.Args())
				if argCount < 0 || argCount > 2 {
					fmt.Println("Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE EXERCISE")
					return
				}

				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
					return
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					if argCount == 0 || argCount == 1 {
						fmt.Println("Are you sure you are logged in? Please login again.")
						c, err = login(configPath)
						if err != nil {
							fmt.Println(err)
							return
						}
					} else {
						c = config.Demo()
					}
				}

				assignments, err := FetchAssignments(c, FetchEndpoint(ctx.Args()))
				if err != nil {
					fmt.Println(err)
					return
				}

				if len(assignments) == 0 {
					noAssignmentMessage := "No assignments found"
					if argCount == 2 {
						fmt.Printf("%s for %s - %s\n", noAssignmentMessage, ctx.Args()[0], ctx.Args()[1])
					} else if argCount == 1 {
						fmt.Printf("%s for %s\n", noAssignmentMessage, ctx.Args()[0])
					} else {
						fmt.Printf("%s\n", noAssignmentMessage)
					}
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(c.Dir, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", c.Dir)
			},
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Save exercism.io api credentials",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				// ignore errors, we're just going to overwrite it anyway
				normalizeConfigFile(configPath)

				_, err := login(config.WithDefaultPath(configPath))
				if err != nil {
					fmt.Println(err)
				}
			},
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     "Clear exercism.io api credentials",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
				}

				logout(config.WithDefaultPath(configPath))
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore completed and current assignments from exercism.io",
			Description: "Restore will pull the latest revisions of exercises that have already been " +
				"submitted. It will *not* overwrite existing files.  If you have made changes " +
				"to a file and have not submitted it, and you're trying to restore the last " +
				"submitted version, first move that file out of the way, then call restore.",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				assignments, err := FetchAssignments(c, FetchEndpoints["restore"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(c.Dir, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", c.Dir)
			},
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     "Submit code to exercism.io on your current assignment",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				if len(ctx.Args()) == 0 {
					fmt.Println("Please enter a file name")
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
					fmt.Println("It looks like this is a test, please enter an example file name.")
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

			},
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     "Delete the last submission",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				err := normalizeConfigFile(configPath)
				if err != nil {
					fmt.Println(err)
				}
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				err = UnsubmitAssignment(c)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("The last submission was successfully deleted.")
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}

func login(path string) (*config.Config, error) {
	c, err := askForConfigInfo()
	if err != nil {
		return nil, err
	}
	err = c.ToFile(path)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Your credentials have been written to %s\n", path)
	fmt.Printf("Your exercism directory can be found at %s\n", c.Dir)
	return c, nil
}

func logout(path string) {
	os.Remove(path)
}

func absolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

func askForConfigInfo() (*config.Config, error) {
	var key, dir string
	delim := "\r\n"

	bio := bufio.NewReader(os.Stdin)

	fmt.Print("Your Exercism API key (found at http://exercism.io/account): ")
	key, err := bio.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Println("What is your exercism exercises project path?")
	fmt.Printf("Press Enter to select the default (%s):\n", config.DefaultAssignmentPath())
	fmt.Print("> ")
	dir, err = bio.ReadString('\n')
	if err != nil {
		return nil, err
	}

	key = strings.TrimRight(key, delim)
	dir = strings.TrimRight(dir, delim)

	if dir == "" {
		dir = config.DefaultAssignmentPath()
	}

	dir = config.ReplaceTilde(dir)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = fmt.Errorf("Error making directory %v: [%v]", dir, err)
		return nil, err
	}

	dir, err = absolutePath(dir)
	if err != nil {
		return nil, err
	}

	return &config.Config{
		APIKey:   key,
		Dir:      dir,
		Hostname: "http://exercism.io",
	}, nil
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

func FetchAssignments(c *config.Config, path string) ([]Assignment, error) {
	url := fmt.Sprintf("%s%s?key=%s", c.Hostname, path, c.APIKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error fetching assignments: [%v]", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		err = fmt.Errorf("Error fetching assignments: [%v]", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error string `json:"error"`
		}
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			err = fmt.Errorf("Error parsing API response: [%v]", err)
			return nil, err
		}

		err = fmt.Errorf("Error fetching assignments. HTTP Status Code: %d\n%s", resp.StatusCode, apiError.Error)
		return nil, err
	}

	var fr struct {
		Assignments []Assignment
	}

	err = json.Unmarshal(body, &fr)
	if err != nil {
		err = fmt.Errorf("Error parsing API response: [%v]", err)
		return nil, err
	}

	return fr.Assignments, nil
}

func UnsubmitAssignment(c *config.Config) error {
	path := "api/v1/user/assignments"

	url := fmt.Sprintf("%s/%s?key=%s", c.Hostname, path, c.APIKey)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error destroying submission: [%v]", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {

		var ur struct {
			Error string
		}

		err = json.Unmarshal(body, &ur)
		if err != nil {
			return err
		}

		err = fmt.Errorf("Status: %d, Error: %v", resp.StatusCode, ur.Error)
		return err
	}

	return nil
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

	req.Header.Set("User-Agent", UserAgent)
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

type Assignment struct {
	Track   string
	Slug    string
	Files   map[string]string
	IsFresh bool `json:"fresh"`
}

func SaveAssignment(dir string, a Assignment) error {
	root := fmt.Sprintf("%s/%s/%s", dir, a.Track, a.Slug)

	for name, text := range a.Files {
		file := fmt.Sprintf("%s/%s", root, name)
		dir := filepath.Dir(file)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("Error making directory %v: [%v]", dir, err)
		}
		if _, err = os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				err = ioutil.WriteFile(file, []byte(text), 0644)
				if err != nil {
					return fmt.Errorf("Error writing file %v: [%v]", name, err)
				}
			}
		}
	}

	fresh := " "
	if a.IsFresh {
		fresh = "*"
	}
	fmt.Println(fresh, a.Track, "-", a.Slug)

	return nil
}

func FetchEndpoint(args []string) string {
	if len(args) == 0 {
		return FetchEndpoints["current"]
	}

	endpoint := FetchEndpoints["exercise"]
	for _, arg := range args {
		endpoint = fmt.Sprintf("%s/%s", endpoint, arg)
	}

	return endpoint
}

func IsTest(filename string) bool {
	for _, ext := range testExtensions {
		if strings.LastIndex(filename, ext) > 0 {
			return true
		}
	}
	return false
}

func normalizeConfigFile(path string) error {
	if path == "" {
		path = config.HomeDir()
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("expected path to be a directory")
	}

	correctPath := filepath.Join(path, config.File)
	legacyPath := filepath.Join(path, config.LegacyFile)

	_, err = os.Stat(correctPath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	_, err = os.Stat(legacyPath)
	if os.IsNotExist(err) {
		return nil
	}

	return os.Rename(legacyPath, correctPath)
}
