package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Unsubmit deletes the most recent submission from the API.
func Unsubmit(ctx *cli.Context) {
	if len(ctx.Args()) > 0 {
		log.Fatal("\nThe unsubmit command does not take any arguments, it deletes the most recent submission.\n\nTo delete a different submission, you'll need to do it from the website.")
	}

	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if !c.IsAuthenticated() {
		log.Fatal(msgPleaseAuthenticate)
	}

	client := api.NewClient(c)
	if err := client.Unsubmit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your most recent submission was successfully deleted.")
}
