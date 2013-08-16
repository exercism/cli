package main

import (
	"exercism"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
)

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = "1.0.0.beta"
	app.Commands = []cli.Command{
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     "Fetch first assignment for each language from exercism.io",
			Action: func(c *cli.Context) {
				config, err := exercism.ConfigFromFile(homeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}
				assignments, err := exercism.FetchAssignments(config.Hostname,
					exercism.FetchEndpoints["demo"], config.ApiKey)
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := exercism.SaveAssignment(config.ExercismDirectory, a)
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
				config, err := exercism.ConfigFromFile(homeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}
				assignments, err := exercism.FetchAssignments(config.Hostname,
					exercism.FetchEndpoints["current"], config.ApiKey)
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := exercism.SaveAssignment(config.ExercismDirectory, a)
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
				usr, err := user.Current()
				if err != nil {
					panic(nil)
				}
				exercism.ConfigToFile(*usr, homeDir(), askForConfigInfo())
			},
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     "Clear exercism.io api credentials",
			Action: func(c *cli.Context) {
				exercism.Logout(homeDir())
			},
		},
		{
			Name:      "peek",
			ShortName: "p",
			Usage:     "Fetch upcoming assignment from exercism.io",
			Action: func(c *cli.Context) {
				config, err := exercism.ConfigFromFile(homeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}
				assignments, err := exercism.FetchAssignments(config.Hostname,
					exercism.FetchEndpoints["next"], config.ApiKey)
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := exercism.SaveAssignment(config.ExercismDirectory, a)
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
				config, err := exercism.ConfigFromFile(homeDir())
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}

				if len(c.Args()) == 0 {
					fmt.Println("Please enter a file name")
					return
				}

				filename := c.Args()[0]

				if exercism.IsTest(filename) {
					fmt.Println("It looks like this is a test, please enter an example file name.")
					return
				}

				submissionPath, err := exercism.SubmitAssignment(config.Hostname, config.ApiKey, filename)
				if err != nil {
					fmt.Printf("There was an issue with your submission: %s\n", err)
					return
				}

				fmt.Printf("For feedback on your submission visit %s%s.",
					config.Hostname, submissionPath)

			},
		},
		{
			Name:      "whoami",
			ShortName: "w",
			Usage:     "Get the github username that you are logged in as",
			Action: func(c *cli.Context) {
				config, err := exercism.ConfigFromFile(homeDir())
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

func homeDir() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.HomeDir
}

func askForConfigInfo() (c exercism.Config) {
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

	if dir == "" {
		dir = currentDir
	}

	return exercism.Config{un, key, dir, "http://exercism.io"}
}
