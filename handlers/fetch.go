package handlers

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Fetch returns exercism problems.
func Fetch(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		fmt.Println(err)
		return
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
		fmt.Println("Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE PROBLEM")
	}

	problems, err := api.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	hw := NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	hw.Report()
	fmt.Println()
}
