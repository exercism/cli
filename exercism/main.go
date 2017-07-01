package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/cmd"
	"github.com/exercism/cli/debug"
	"github.com/urfave/cli"
)

const (
	// Version is the current release of the command-line app.
	// We try to follow Semantic Versioning (http://semver.org),
	// but with the http://exercism.io app being a prototype, a
	// lot of things get out of hand.
	Version = "2.4.1"

	descConfigure = "Writes config values to a JSON file."
	descDebug     = "Outputs useful debug information."
	descDownload  = "Downloads a solution given the ID of the latest iteration."
	descFetch     = "Fetches the next unsubmitted problem in each track."
	descList      = "Lists the available problems for a language track, given its ID."
	descOpen      = "Opens exercism.io to your most recent iteration of a problem given the track ID and problem slug."
	descRestore   = "Downloads the most recent iteration for each of your solutions on exercism.io."
	descSkip      = "Skips a problem given a track ID and problem slug."
	descStatus    = "Fetches information about your progress with a given language track."
	descSubmit    = "Submits a new iteration to a problem on exercism.io."
	descTracks    = "Lists the available language tracks."
	descUpgrade   = "Upgrades the CLI to the latest released version."

	descLongDownload = "The submission ID is the last part of the URL when looking at a solution on exercism.io."
	descLongRestore  = "Restore will pull the latest revisions of exercises that have already been submitted. It will *not* overwrite existing files. If you have made changes to a file and have not submitted it, and you're trying to restore the last submitted version, first move that file out of the way, then call restore."
)

func main() {
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)

	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = Version
	app.Before = func(ctx *cli.Context) error {
		debug.Verbose = ctx.GlobalBool("verbose")
		debug.Println("verbose logging enabled")

		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "path to config file",
			EnvVar: "EXERCISM_CONFIG_FILE,XDG_CONFIG_HOME",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "turn on verbose logging",
		},
	}
	app.Commands = []cli.Command{
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
					Usage: "exercism.io API key (see http://exercism.io/account/key)",
				},
				cli.StringFlag{
					Name:  "api, a",
					Usage: "exercism xapi host",
				},
				cli.BoolFlag{
					Name:  "silent, s",
					Usage: "Obfuscates configuration options from output",
				},
			},
			Action: cmd.Configure,
		},
		{
			Name:  "debug",
			Usage: descDebug,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "full-api-key",
					Usage: "Displays the full API key without obfuscating it",
				},
			},
			Action: cmd.Debug,
		},
		{
			Name:        "download",
			ShortName:   "dl",
			Usage:       descDownload,
			Description: descLongDownload,
			Action:      cmd.Download,
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     descFetch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all",
					Usage: "fetch all exercises for a given track",
				},
			},
			Action: cmd.Fetch,
		},
		{
			Name:      "list",
			ShortName: "li",
			Usage:     descList,
			Action:    cmd.List,
		},
		{
			Name:      "open",
			ShortName: "op",
			Usage:     descOpen,
			Action:    cmd.Open,
		},
		{
			Name:        "restore",
			ShortName:   "r",
			Usage:       descRestore,
			Description: descLongRestore,
			Action:      cmd.Restore,
		},
		{
			Name:   "skip",
			Usage:  descSkip,
			Action: cmd.Skip,
		},
		{
			Name:      "status",
			ShortName: "st",
			Usage:     descStatus,
			Action:    cmd.Status,
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     descSubmit,
			Action:    cmd.Submit,
			ArgsUsage: "file1 [file2, etc...]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "test",
					Usage: "allow submission of test files",
				},
				cli.StringFlag{
					Name:  "comment, m",
					Usage: "includes a comment with the submission",
				},
			},
		},
		{
			Name:      "tracks",
			ShortName: "t",
			Usage:     descTracks,
			Action:    cmd.Tracks,
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     "REMOVED",
			Action: func(*cli.Context) {
				fmt.Println("For security reasons, this command is no longer in use.\nYou can delete iterations in the web interface.")
			},
		},
		{
			Name:   "upgrade",
			Usage:  descUpgrade,
			Action: cmd.Upgrade,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
