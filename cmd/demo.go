package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"../api"
	"../config"
	"../user"
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
