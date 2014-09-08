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

	hw := NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	hw.Report(HWAll)

	fmt.Println("Next step: choose a language, read the README, and make the test suite pass.")
}
