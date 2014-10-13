package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
)

// Demo returns one problem for each active track.
func Demo(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	problems, err := api.Demo(c)
	if err != nil {
		log.Fatal(err)
	}

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}

	hw.Report(user.HWAll)

	fmt.Println("Next step: choose a language, read the README, and make the test suite pass.")
}
