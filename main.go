package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/handlers"
)

const (
	// Version is the current release of the command-line app.
	// We try to follow Semantic Versioning (http://semver.org),
	// but with the http://exercism.io app being a prototype, a
	// lot of things get out of hand.
	Version = "1.7.4"

	descDebug     = "Outputs useful debug information."
	descConfigure = "Writes config values to a JSON file."
	descDemo      = "Fetches a demo problem for each language track on exercism.io."
	descFetch     = "Fetches your current problems on exercism.io, as well as the next unstarted problem in each language."
	descRestore   = "Restores completed and current problems on from exercism.io, along with your most recent iteration for each."
	descSubmit    = "Submits a new iteration to a problem on exercism.io."
	descUnsubmit  = "Deletes the most recently submitted iteration."
	descLogin     = "DEPRECATED: Interactively saves exercism.io api credentials."
	descLogout    = "DEPRECATED: Clear exercism.io api credentials"

	descLongRestore = "Restore will pull the latest revisions of exercises that have already been submitted. It will *not* overwrite existing files. If you have made changes to a file and have not submitted it, and you're trying to restore the last submitted version, first move that file out of the way, then call restore."
)

func main() {
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)

	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "debug",
			Usage:  descDebug,
			Action: handlers.Debug,
		},
		{
			Name:  "configure",
			Usage: descConfigure,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "path to exercises directory",
				},
				cli.StringFlag{
					Name:  "host, u",
					Usage: "exercism api host",
				},
				cli.StringFlag{
					Name:  "key, k",
					Usage: "exercism.io API key (see http://exercism.io/account)",
				},
			},
			Action: handlers.Configure,
		},
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     descDemo,
			Action:    handlers.Demo,
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     descFetch,
			Action:    handlers.Fetch,
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     descLogin,
			Action:    handlers.Login,
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     descLogout,
			Action:    handlers.Logout,
		},
		{
			Name:        "restore",
			ShortName:   "r",
			Usage:       descRestore,
			Description: descLongRestore,
			Action:      handlers.Restore,
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     descSubmit,
			Action:    handlers.Submit,
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     descUnsubmit,
			Action:    handlers.Unsubmit,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
