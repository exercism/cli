package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func Home(ctx *cli.Context) {
	path, err := config.Path(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.FromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Your exercism directory can be found at %s\n", c.Dir)
}
