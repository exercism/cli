package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"../api"
	"../config"
	"../user"
)

// Tracks lists available tracks.
func Tracks(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	tracks, err := api.Tracks(fmt.Sprintf("%s/tracks", c.XAPI))
	if err != nil {
		log.Fatal(err)
	}

	curr := user.NewCurriculum(tracks)
	fmt.Println("\nActive language tracks:")
	curr.Report(user.TrackActive)
	fmt.Println("\nInactive language tracks:")
	curr.Report(user.TrackInactive)

	// TODO: implement `list` command to list problems in a track
	msg := `
Related commands:
    exercism fetch (see 'exercism help fetch')
	`
	fmt.Println(msg)
}
