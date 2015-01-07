package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"../api"
	"../config"
	"../user"
)

// Fetch returns exercism problems.
func Fetch(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	args := ctx.Args()
	var url string
	switch len(args) {
	case 0:
		url = fmt.Sprintf("%s/%s?key=%s", c.XAPI, "v2/exercises", c.APIKey)
	case 1:
		url = fmt.Sprintf("%s/%s/%s?key=%s", c.XAPI, "v2/exercises", args[0], c.APIKey)
	case 2:
		url = fmt.Sprintf("%s/%s/%s/%s", c.XAPI, "v2/exercises", args[0], args[1])
	default:
		msg := "Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE PROBLEM"
		log.Fatal(msg)
	}

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
