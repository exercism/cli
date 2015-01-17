package cmd

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
	defer fmt.Printf("\nIf you are having trouble and need to file a GitHub issue (https://github.com/exercism/exercism.io/issues) please include this information (except your API key. Keep that private).\n")

	fmt.Printf("\n**** Debug Information ****\n")
	fmt.Printf("Exercism CLI Version: %s\n", ctx.App.Version)
	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	dir, err := config.Home()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Home Dir: %s\n", dir)

	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	configured := true
	if _, err = os.Stat(c.File); err != nil {
		if os.IsNotExist(err) {
			configured = false
		} else {
			log.Fatal(err)
		}
	}

	if configured {
		fmt.Printf("Config file: %s\n", c.File)
		fmt.Printf("API Key: %s\n", c.APIKey)
	} else {
		fmt.Println("Config file: <not configured>")
		fmt.Println("API Key: <not configured>")
	}
	fmt.Printf("API: %s\n", c.API)
	fmt.Printf("Exercises Directory: %s\n", c.Dir)
}
