package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/rpt"
)

// Restore returns a user's solved problems.
func Restore(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/api/v1/iterations/%s/restore", c.API, c.APIKey)

	problems, err := api.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	hw := rpt.NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		log.Fatal(err)
	}
	hw.Summarize()
}
