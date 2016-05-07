package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Open uses the given track and problem and opens it in the browser.
func Open(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: exercism open TRACK_ID PROBLEM")
		os.Exit(1)
	}

	trackID := args[0]
	slug := args[1]
	submission, err := client.SubmissionURL(trackID, slug)
	if err != nil {
		log.Fatal(err)
	}

	url := submission.URL
	// Escape characters are not allowed by cmd/bash.
	switch runtime.GOOS {
	case "windows":
		url = strings.Replace(url, "&", `^&`, -1)
	default:
		url = strings.Replace(url, "&", `\&`, -1)
	}

	// The command to open the browser is OS-dependent.
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "freebsd", "linux", "netbsd", "openbsd":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
