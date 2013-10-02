package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = VERSION
	app.Commands = []cli.Command{
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     "Fetch first assignment for each language from exercism.io",
			Action: func(c *cli.Context) {
				config, err := ConfigFromFile(HomeDir())
				if err != nil {
					demoDir, err := DemoDirectory()
					if err != nil {
						fmt.Println(err)
						return
					}
					config = Config{
						Hostname:          "http://exercism.io",
						ApiKey:            "",
						ExercismDirectory: demoDir,
					}
				}
				assignments, err := FetchAssignments(config,
					FetchEndpoints["demo"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}
			},
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     "Fetch current assignment from exercism.io",
			Action: func(c *cli.Context) {
				config, err := ConfigFromFile(HomeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}
				assignments, err := FetchAssignments(config,
					FetchEndpoints["current"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}
			},
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Save exercism.io api credentials",
			Action: func(c *cli.Context) {
				ConfigToFile(HomeDir(), askForConfigInfo())
			},
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     "Clear exercism.io api credentials",
			Action: func(c *cli.Context) {
				Logout(HomeDir())
			},
		},
		{
			Name:      "peek",
			ShortName: "p",
			Usage:     "Fetch upcoming assignment from exercism.io",
			Action: func(c *cli.Context) {
				config, err := ConfigFromFile(HomeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}
				assignments, err := FetchAssignments(config,
					FetchEndpoints["next"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}
			},
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     "Submit code to exercism.io on your current assignment",
			Action: func(c *cli.Context) {
				config, err := ConfigFromFile(HomeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}

				if len(c.Args()) == 0 {
					fmt.Println("Please enter a file name")
					return
				}

				filename := c.Args()[0]

				// Make filename relative to config.ExercismDirectory.
				absPath, err := absolutePath(filename)
				if err != nil {
					fmt.Printf("Couldn't find %v: %v\n", filename, err)
					return
				}
				exDir := config.ExercismDirectory + string(filepath.Separator)
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

				response, err := SubmitAssignment(config, filename, code)
				if err != nil {
					fmt.Printf("There was an issue with your submission: %v\n", err)
					return
				}

				fmt.Printf("For feedback on your submission visit %s%s.\n",
					config.Hostname, response.SubmissionPath)

			},
		},
		{
			Name:      "whoami",
			ShortName: "w",
			Usage:     "Get the github username that you are logged in as",
			Action: func(c *cli.Context) {
				config, err := ConfigFromFile(HomeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}

				fmt.Println(config.GithubUsername)
			},
		},
	}
	app.Run(os.Args)
}

func askForConfigInfo() (c Config) {
	var un, key, dir string

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Print("Your GitHub username: ")
	_, err = fmt.Scanln(&un)
	if err != nil {
		panic(err)
	}

	fmt.Print("Your exercism.io API key: ")
	_, err = fmt.Scanln(&key)
	if err != nil {
		panic(err)
	}

	fmt.Println("What is your exercism exercises project path?")
	fmt.Printf("Press Enter to select the default (%s):\n", currentDir)
	fmt.Print("> ")
	_, err = fmt.Scanln(&dir)
	if err != nil && err.Error() != "unexpected newline" {
		panic(err)
	}
	dir, err = absolutePath(dir)
	if err != nil {
		panic(err)
	}

	if dir == "" {
		dir = currentDir
	}

	return Config{un, key, ReplaceTilde(dir), "http://exercism.io"}
}

func absolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}
