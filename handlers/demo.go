package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
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

	hw := NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		log.Fatal(err)
	}

	hw.Report(HWAll)

	fmt.Println("Next step: choose a language, read the README, and make the test suite pass.")
}
