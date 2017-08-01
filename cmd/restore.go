package cmd

import (
	"log"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
	app "github.com/urfave/cli"
)

// Restore returns a user's solved problems.
func Restore(ctx *app.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient(c)

	problems, err := client.Restore()
	if err != nil {
		log.Fatal(err)
	}

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}
	hw.Summarize(user.HWNotSubmitted)

	return nil
}
