package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

const msgExplainFetch = "In order to fetch a specific assignment, call the fetch command with a specific assignment.\n\nexercism fetch ruby matrix"

// List returns the full list of assignments for a given track.
func List(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 1 {
		msg := "Usage: exercism list LANGUAGE"
		log.Fatal(msg)
	}

	language := args[0]
	client := api.NewClient(c)
	problems, err := client.List(language)
	if err != nil {
		if err == api.ErrUnknownLanguage {
			log.Fatalf("The requested language '%s' is unknown", language)
		}
		log.Fatal(err)
	}

	for _, p := range problems {
		fmt.Printf("%s\n", p)
	}
	fmt.Printf("\n%s\n\n", msgExplainFetch)
}
