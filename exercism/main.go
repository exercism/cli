package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/cmd"
)

const (
	// Version is the current release of the command-line app.
	// We try to follow Semantic Versioning (http://semver.org),
	// but with the http://exercism.io app being a prototype, a
	// lot of things get out of hand.
	Version = "1.8.2"

	descDebug     = "Outputs useful debug information."
	descConfigure = "Writes config values to a JSON file."
	descDemo      = "Fetches a demo problem for each language track on exercism.io."
	descFetch     = "Fetches your current problems on exercism.io, as well as the next unstarted problem in each language."
	descRestore   = "Restores completed and current problems on from exercism.io, along with your most recent iteration for each."
	descSubmit    = "Submits a new iteration to a problem on exercism.io."
	descUnsubmit  = "Deletes the most recently submitted iteration."
	descTracks    = "List the available language tracks"
	descLogin     = "DEPRECATED: Interactively saves exercism.io api credentials."
	descLogout    = "DEPRECATED: Clear exercism.io api credentials"

	descLongRestore = "Restore will pull the latest revisions of exercises that have already been submitted. It will *not* overwrite existing files. If you have made changes to a file and have not submitted it, and you're trying to restore the last submitted version, first move that file out of the way, then call restore."
	descDownload 	= "Downloads and saves a specified submission into the local system"
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
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "turn on verbose logging",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "debug",
			Usage:  descDebug,
			Action: cmd.Debug,
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
			Action: cmd.Configure,
		},
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     descDemo,
			Action:    cmd.Demo,
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     descFetch,
			Action:    cmd.Fetch,
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     descLogin,
			Action:    cmd.Login,
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     descLogout,
			Action:    cmd.Logout,
		},
		{
			Name:        "restore",
			ShortName:   "r",
			Usage:       descRestore,
			Description: descLongRestore,
			Action:      cmd.Restore,
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     descSubmit,
			Action:    cmd.Submit,
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     descUnsubmit,
			Action:    cmd.Unsubmit,
		},
		{
			Name:      "tracks",
			ShortName: "t",
			Usage:     descTracks,
			Action:    cmd.Tracks,
		},
		{
			Name:      "download",
			ShortName: "dl",
			Usage:     descDownload,
			Action:    cmd.Download,
		},
	}
	if err := app.Run(os.Args); err != nil {

		log.Fatal(err)
	}
}
