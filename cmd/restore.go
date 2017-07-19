package cmd

import (
	"fmt"
	"log"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
	"github.com/urfave/cli"
)

// Restore returns a user's solved problems.
func Restore(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient(c)

	var problems []*api.Problem

	switch {
	case len(ctx.Args()) == 0 && ctx.Bool("force"):
		fmt.Printf("You are trying to restore all exercises at once, this can take a while, please stay patient")
		if problems, err = client.Restore(); err != nil {
			log.Fatal(err)
		}
	case len(ctx.Args()) == 0:
		log.Fatalf("Restoring everything could take quite a while, please use `--force` if you are sure.")
	}

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}
	hw.Summarize(user.HWNotSubmitted)

	return nil
}
