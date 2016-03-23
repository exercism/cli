package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/paths"
)

// Submit posts an iteration to the API.
func Submit(ctx *cli.Context) {
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

		file, err = filepath.EvalSymlinks(file)
		if err != nil {
			log.Fatal(err)
		}

		if ctx.GlobalBool("verbose") {
			log.Printf("eval symlinks (file): %s", file)
		}

		files = append(files, file)
	}

	iteration, err := api.NewIteration(dir, files)
	if err != nil {
		log.Fatalf("Unable to submit - %s", err)
	}
	iteration.Key = c.APIKey
	iteration.Comment = ctx.String("comment")

	client := api.NewClient(c)
	submission, err := client.Submit(iteration)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s - %s\n%s\n\n", submission.Language, submission.Name, submission.URL)
}
