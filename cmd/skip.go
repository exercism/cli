package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
	"github.com/urfave/cli"
)

// Skip allows a user to skip a specific exercise.
func Skip(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: exercism skip TRACK_ID EXERCISE\n")
		os.Exit(1)
	}

	var (
		trackID = args[0]
		slug    = args[1]
	)

	client := api.NewClient(c)
	if err := client.Skip(trackID, slug); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Exercise %q in %q has been skipped.\n", slug, trackID)

	return nil
}
