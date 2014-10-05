package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
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

	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if !c.IsAuthenticated() {
		log.Fatal(msgPleaseAuthenticate)
	}

	filename := ctx.Args()[0]

	if isTest(filename) {
		log.Fatal("Please submit the solution, not the test file.")
	}

	file, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}
	file, err = filepath.EvalSymlinks(file)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := filepath.EvalSymlinks(c.Dir)
	if err != nil {
		log.Fatal(err)
	}

	code, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot read the contents of %s - %s\n", filename, err)
	}

	url := fmt.Sprintf("%s/api/v1/user/assignments", c.API)

	iteration := &api.Iteration{
		Key:  c.APIKey,
		Code: string(code),
		File: file,
		Dir:  dir,
	}

	err = iteration.Identify()
	if err != nil {
		msg := `Please leave the solution within the problem directory that was created by 'exercism fetch'`
		log.Fatalf("Cannot submit - %s.\n\n%s", err, msg)
	}

	submission, err := api.Submit(url, iteration)
	if err != nil {
		log.Fatal(err)
	}

	msg := "Submitted %s in %s. Your submission can be found online at %s\n"
	fmt.Printf(msg, submission.Name, submission.Language, submission.URL)
}
