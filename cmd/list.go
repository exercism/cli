package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

const msgExplainFetch = "In order to fetch a specific assignment, call the fetch command with a specific assignment.\n\nexercism fetch %s %s\n\n"

// List returns the full list of assignments for a given track.
func List(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: exercism list TRACK_ID\n")
		os.Exit(1)
	}

	trackID := args[0]
	client := api.NewClient(c)
	problems, err := client.List(trackID)
	if err != nil {
		if err == api.ErrUnknownTrack {
			log.Fatalf("There is no track with ID '%s'.", trackID)
		}
		log.Fatal(err)
	}

	for _, p := range problems {
		fmt.Printf("%s\n", p)
	}
	fmt.Println()
	fmt.Printf(msgExplainFetch, trackID, problems[0])

	return nil
}
