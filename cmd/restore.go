package cmd

import (
	"log"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
	"github.com/robphoenix/cli/user"
	"github.com/urfave/cli"
)

// Restore returns a user's solved exercises.
func Restore(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient(c)

	exercises, err := client.Restore()
	if err != nil {
		log.Fatal(err)
	}

	hw := user.NewHomework(exercises, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}
	hw.Summarize(user.HWNotSubmitted)

	return nil
}
