package handlers

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

// Debug provides information about the user's environment and configuration.
func Debug(ctx *cli.Context) {
	defer fmt.Printf("\nIf you are having any issues, please contact kytrinyx@exercism.io with this information.\n")

	bail := func(err error) {
		if err != nil {
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
	configured := true
	if _, err = os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			configured = false
		} else {
			bail(err)
		}
	}

	c, err := config.Read(file)
	bail(err)

	if configured {
		fmt.Printf("Config file: %s\n", c.File())
		fmt.Printf("API Key: %s\n", c.APIKey)
	} else {
		fmt.Println("Config file: <not configured>")
		fmt.Println("API Key: <not configured>")
	}
	fmt.Printf("API: %s\n", c.Hostname)
	fmt.Printf("Exercises Directory: %s\n", c.Dir)
}
