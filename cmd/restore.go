package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"../api"
	"../config"
	"../user"
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

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}
	hw.Summarize()
}
