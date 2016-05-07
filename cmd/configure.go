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
func Configure(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	key := ctx.String("key")
	host := ctx.String("host")
	dir := ctx.String("dir")
	api := ctx.String("api")

	if err := c.Update(key, host, dir, api); err != nil {
		log.Fatalf("Error updating your configuration %s\n", err)
	}

	if err := os.MkdirAll(c.Dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating exercism directory %s\n", err)
	}

	if err := c.Write(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nConfiguration written to %s\n\n", c.File)
	fmt.Printf("  --key=%s\n", c.APIKey)
	fmt.Printf("  --dir=%s\n", c.Dir)
	fmt.Printf("  --host=%s\n", c.API)
	fmt.Printf("  --api=%s\n\n", c.XAPI)

	return nil
}
