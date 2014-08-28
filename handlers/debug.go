package handlers

import (
	"fmt"
	"log"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func Debug(ctx *cli.Context) {
	bail := func(err error) {
		if err != nil {
			fmt.Printf("\nIf you are having any issues, please contact kytrinyx@exercism.io with this information.\n")
			log.Fatal(err)
		}
	}

	fmt.Printf("\n**** Debug Information ****\n")
	fmt.Printf("Exercism CLI Version: %s\n", ctx.App.Version)
	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	dir, err := config.Home()
	bail(err)
	fmt.Printf("Home Dir: %s\n", dir)

	file, err := config.FilePath(ctx.GlobalString("config"))
	bail(err)

	c, err := config.Read(file)
	bail(err)

	fmt.Printf("Config file: %s\n", c.File())
	fmt.Printf("API: %s\n", c.Hostname)
	fmt.Printf("API Key: %s\n", c.APIKey)
	fmt.Printf("Exercises Directory: %s\n", c.Dir)
	fmt.Println()
}
