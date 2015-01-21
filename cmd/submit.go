package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Submit posts an iteration to the api
func Submit(ctx *cli.Context) {
	if len(ctx.Args()) == 0 {
		log.Fatal("Please enter a file name")
	}

	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("debug") {
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

	filename := ctx.Args()[0]

	if ctx.GlobalBool("debug") {
		log.Printf("file name: %s", filename)
	}

	if isTest(filename) {
		log.Fatal("Please submit the solution, not the test file.")
	}

	file, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("debug") {
		log.Printf("absolute path: %s", file)
	}

	file, err = filepath.EvalSymlinks(file)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("debug") {
		log.Printf("eval symlinks (file): %s", file)
	}

	dir, err := filepath.EvalSymlinks(c.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.GlobalBool("debug") {
		log.Printf("eval symlinks (dir): %s", dir)
	}

	code, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot read the contents of %s - %s\n", filename, err)
	}

	iteration := &api.Iteration{
		Key:  c.APIKey,
		Code: string(code),
		File: file,
		Dir:  dir,
	}

	if err := iteration.Identify(); err != nil {
		msg := `Please leave the solution within the problem directory that was created by 'exercism fetch'`
		log.Fatalf("Cannot submit - %s.\n\n%s", err, msg)
	}

	client := api.NewClient(c)
	submission, err := client.Submit(iteration)
	if err != nil {
		log.Fatal(err)
	}

	msg := `
Submitted %s in %s.
Your submission can be found online at %s

To get the next exercise, run "exercism fetch" again.
`
	fmt.Printf(msg, submission.Name, submission.Language, submission.URL)
}
