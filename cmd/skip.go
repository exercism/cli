package cmd

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Skip command allows a user to skip a specific problem
func Skip(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	args := ctx.Args()

	if len(args) != 2 {
		msg := "Usage: exercism skip LANGUAGE SLUG"
		log.Fatal(msg)
	}

	var (
		language = args[0]
		slug     = args[1]
	)

	client := api.NewClient(c)
	if err := client.Skip(language, slug); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Exercise %q in %q has been skipped.\n", slug, language)
}
