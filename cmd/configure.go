package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

// Configure stores settings in a JSON file.
// If a setting is not passed as an argument, default
// values are used.
func Configure(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	key := ctx.String("key")
	host := ctx.String("host")
	dir := ctx.String("dir")
	api := ctx.String("api")
	c.Update(key, host, dir, api)

	if err := os.MkdirAll(c.Dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating exercism directory %s\n", err)
	}

	if err := c.Write(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The configuration has been written to %s\n", c.File)
	fmt.Printf("Your exercism directory can be found at %s\n", c.Dir)
}
