package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Status is a command that allows a user to view their progress in a given
// language track.
func Status(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 1 {
		log.Fatal("Usage: exercism status TRACK_ID")
	}

	client := api.NewClient(c)
	trackID := args[0]
	status, err := client.Status(trackID)
	if err != nil {
		if err == api.ErrUnknownTrack {
			log.Fatalf("There is no track with ID '%s'.", trackID)
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println(status)
}
