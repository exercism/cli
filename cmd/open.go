package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/browser"
	"github.com/exercism/cli/config"
	app "github.com/urfave/cli"
)

// Open opens the user's latest iteration of the exercise on the given track.
func Open(ctx *app.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: exercism open TRACK_ID PROBLEM\n")
		os.Exit(1)
	}

	trackID := args[0]
	slug := args[1]
	submission, err := client.SubmissionURL(trackID, slug)
	if err != nil {
		return err
	}

	return browser.Open(submission.URL)
}
