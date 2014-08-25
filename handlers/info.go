package handlers

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func Info(ctx *cli.Context) {
	path, err := config.Path(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.FromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("API Key:", c.APIKey)
	fmt.Println("Exercises Directory:", c.Dir)
	fmt.Println("Config file:", c.File())
	fmt.Println("API:", c.Hostname)
}
