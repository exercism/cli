package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
	"github.com/robphoenix/cli/paths"
	"github.com/urfave/cli"
)

// Submit posts an iteration to the API.
func Submit(ctx *cli.Context) error {
	if len(ctx.Args()) == 0 {
		log.Fatal("Please enter a file name")
	}

	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("verbose") {
		log.Printf("Exercises dir: %s", c.Dir)
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("Unable to get current working directory - %s", err)
		} else {
			log.Printf("Current dir: %s", dir)
		}
	}

	if !c.IsAuthenticated() {
		log.Fatal(msgPleaseAuthenticate)
	}

	dir, err := filepath.EvalSymlinks(c.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("verbose") {
		log.Printf("eval symlinks (dir): %s", dir)
	}

	files := []string{}
	for _, filename := range ctx.Args() {
		if ctx.GlobalBool("verbose") {
			log.Printf("file name: %s", filename)
		}

		if isTest(filename) && !ctx.Bool("test") {
			log.Fatal("You're trying to submit a test file. If this is really what " +
				"you want, please pass the --test flag to exercism submit.")
		}

		if isREADME(filename) {
			log.Fatal("You cannot submit the README as a solution.")
		}

		if paths.IsDir(filename) {
			log.Fatal("Please specify each file that should be submitted, e.g. `exercism submit file1 file2 file3`.")
		}

		file, err := filepath.Abs(filename)
		if err != nil {
			log.Fatal(err)
		}

		if ctx.GlobalBool("verbose") {
			log.Printf("absolute path: %s", file)
		}
		files = append(files, file)
	}

	iteration, err := api.NewIteration(dir, files)
	if err != nil {
		log.Fatalf("unable to submit - %s", err)
	}
	iteration.Key = c.APIKey
	iteration.Comment = ctx.String("comment")

	client := api.NewClient(c)
	submission, err := client.Submit(iteration)
	if err != nil {
		log.Fatal(err)
	}

	solutionURL, _ := url.Parse(c.API)
	solutionURL.Path += fmt.Sprintf("tracks/%s/exercises/%s", iteration.TrackID, iteration.Exercise)
	fmt.Printf("Your %s solution for %s has been submitted. View it here:\n%s\n\n", submission.Language, submission.Name, submission.URL)
	fmt.Printf("See related solutions and get involved here:\n%s\n\n", solutionURL)

	return nil
}
