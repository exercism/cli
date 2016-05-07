package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Skip allows a user to skip a specific problem.
func Skip(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: exercism skip TRACK_ID PROBLEM")
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
