package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"../config"
)

// Logout deletes the config file.
// Delete when nobody is using 1.6.x anymore.
func Logout(ctx *cli.Context) {
	msg := `
	*******************************************************************
	DEPRECATED!

	In the future use the 'exercism configure' command to reconfigure:

		exercism configure --key YOUR_API_KEY

	Or delete the config file yourself.

	*******************************************************************

	`
	fmt.Printf(msg)

	file, err := config.FilePath(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleting config file at %s\n\n", file)
	os.Remove(file)
}
