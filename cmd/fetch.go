package cmd

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
)

// Fetch returns exercism problems.
func Fetch(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	problems, err := client.Fetch(ctx.Args())
	if err != nil {
		log.Fatal(err)
	}

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}

	hw.Summarize()
}
