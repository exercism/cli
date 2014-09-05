package handlers

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Demo returns one problem for each active track.
func Demo(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		fmt.Println(err)
		return
	}

	problems, err := api.Demo(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, problem := range problems {
		err := problem.Save(c.Dir)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	NewHomework(problems, c).Report()

	fmt.Println()
	fmt.Println("Next step: choose a language, read the README, and make the test suite pass.")
}
