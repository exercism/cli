package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
)

// Tracks lists available tracks.
func Tracks(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	tracks, err := client.Tracks()
	if err != nil {
		log.Fatal(err)
	}

	curr := user.NewCurriculum(tracks)
	fmt.Println("\nActive language tracks:")
	curr.Report(user.TrackActive)
	fmt.Println("\nInactive language tracks:")
	curr.Report(user.TrackInactive)

	msg := `
Related commands:
    exercism fetch (see 'exercism help fetch')
    exercism list (see 'exercism help list')
	`
	fmt.Println(msg)

	return nil
}
