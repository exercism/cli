package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func Configure(ctx *cli.Context) {
	key := ctx.String("key")
	host := ctx.String("host")
	dir := ctx.String("dir")
	c, err := config.New(key, host, dir)
	if err != nil {
		log.Fatal(err)
	}

	c.SavePath(ctx.GlobalString("config"))

	err = c.Write()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The configuration has been written to %s\n", c.File())
	fmt.Printf("Your exercism directory can be found at %s\n", c.Dir)
}
