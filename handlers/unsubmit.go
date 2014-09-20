package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
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
	err = api.Unsubmit(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Your most recent submission was successfully deleted.")
}
