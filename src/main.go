package main

import "os"
import "github.com/codegangsta/cli"

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.Run(os.Args)
}
