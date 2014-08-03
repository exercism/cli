package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", config.Filename(config.HomeDir()), "path to config file"},
	}
	app.Commands = []cli.Command{
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
				c, err := config.FromFile(configPath)
				if err != nil {
					c, err = config.Demo()
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				assignments, err := FetchAssignments(c, FetchEndpoints["demo"])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, a := range assignments {
					err := SaveAssignment(c.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				msg := "\nThe demo exercises have been written to %s, in subdirectories by language.\n\nTo try an exercise, change directory to a language/exercise, read the README and run the tests.\n\n"
				fmt.Printf(msg, c.ExercismDirectory)
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
						c, err = config.Demo()
						if err != nil {
							fmt.Println(err)
							return
						}
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
					err := SaveAssignment(c.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", c.ExercismDirectory)
			},
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Save exercism.io api credentials",
			Action: func(ctx *cli.Context) {
				_, err := login(ctx.GlobalString("config"))
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
				logout(ctx.GlobalString("config"))
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
					err := SaveAssignment(c.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", c.ExercismDirectory)
			},
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     "Submit code to exercism.io on your current assignment",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
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

				// Make filename relative to config.ExercismDirectory.
				absPath, err := absolutePath(filename)
				if err != nil {
					fmt.Printf("Couldn't find %v: %v\n", filename, err)
					return
				}
				exDir := c.ExercismDirectory + string(filepath.Separator)
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
					c.Hostname, "/submissions/", response.Id)

			},
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     "Delete the last submission",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				response, err := UnsubmitAssignment(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				if response != "" {
					return
				}

				fmt.Println("The last submission was successfully deleted.")
			},
		},
		{
			Name:      "whoami",
			ShortName: "w",
			Usage:     "Get the github username that you are logged in as",
			Action: func(ctx *cli.Context) {
				configPath := ctx.GlobalString("config")
				c, err := config.FromFile(configPath)
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					c, err = login(configPath)
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				fmt.Println(c.GithubUsername)
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}
