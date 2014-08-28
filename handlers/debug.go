package handlers

import (
	"fmt"
	"log"
	"os/user"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

func Debug(ctx *cli.Context) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Debug Information\n")
	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Home Dir: %s\n", usr.HomeDir)
	fmt.Printf("Version: %s\n", ctx.App.Version)

	file, err := config.FilePath(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	c, err := config.Read(file)
	if err == nil {
		fmt.Printf("\nExercism Configuration\n")
		fmt.Printf("API Key: %s\n", c.APIKey)
		fmt.Printf("Exercises Directory: %s\n", c.Dir)
		fmt.Printf("Config file: %s\n", c.File())
		fmt.Printf("API: %s\n", c.Hostname)
	}

	fmt.Printf("\nIf you are having any issues, please contact kytrinyx@exercism.io with this information.\n")
}
