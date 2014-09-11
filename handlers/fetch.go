package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
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
		url = fmt.Sprintf("%s/%s?key=%s", c.ProblemsHost, "v2/exercises", c.APIKey)
	case 1:
		url = fmt.Sprintf("%s/%s/%s?key=%s", c.ProblemsHost, "v2/exercises", args[0], c.APIKey)
	case 2:
		url = fmt.Sprintf("%s/%s/%s/%s", c.ProblemsHost, "v2/exercises", args[0], args[1])
	default:
		msg := "Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE PROBLEM"
		log.Fatal(msg)
	}

	problems, err := api.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	hw := NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		log.Fatal(err)
	}

	hw.Summarize()
}
