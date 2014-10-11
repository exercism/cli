package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

// Configure stores settings in a JSON file.
// If a setting is not passed as an argument, default
// values are used.
func Configure(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	c.SavePath(ctx.GlobalString("config"))

	key := ctx.String("key")
	host := ctx.String("host")
	dir := ctx.String("dir")
	c.Update(key, host, dir)

	err = c.Write()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The configuration has been written to %s\n", c.File())
	fmt.Printf("Your exercism directory can be found at %s\n", c.Dir)
}
