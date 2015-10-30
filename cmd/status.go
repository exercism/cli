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
		log.Fatal("Usage: exercism status LANGUAGE")
	}

	client := api.NewClient(c)
	status, err := client.Status(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(status)
}
