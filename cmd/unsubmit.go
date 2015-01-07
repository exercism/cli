package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"../api"
	"../config"
)

// Unsubmit deletes an iteration from the api.
// If no iteration is specified, the most recent iteration
// is deleted.
func Unsubmit(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	if !c.IsAuthenticated() {
		log.Fatal(msgPleaseAuthenticate)
	}

	url := fmt.Sprintf("%s/api/v1/user/assignments?key=%s", c.API, c.APIKey)
	if err := api.Unsubmit(url); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Your most recent submission was successfully deleted.")
}
