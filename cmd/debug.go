package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
)

// Debug provides information about the user's environment and configuration.
func Debug(ctx *cli.Context) {
	defer fmt.Printf("\nIf you are having trouble and need to file a GitHub issue (https://github.com/exercism/exercism.io/issues) please include this information (except your API key. Keep that private).\n")

	client := http.Client{Timeout: 5 * time.Second}

	fmt.Printf("\n**** Debug Information ****\n")
	fmt.Printf("Exercism CLI Version: %s\n", ctx.App.Version)

	rel, err := checkLatestRelease(client)
	if err != nil {
		log.Println(err)
	} else {
		if rel.Version() != ctx.App.Version {
			defer fmt.Printf("\nA newer version of the CLI (%s) can be downloaded here: %s\n", rel.TagName, rel.Location)
		}
		fmt.Printf("Exercism CLI Latest Release: %s\n", rel.Version())
	}

	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	dir, err := config.Home()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Home Dir: %s\n", dir)

	c, err := config.New(ctx.GlobalString("config"))
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

	fmt.Printf("API: %s [%s]\n", c.API, pingUrl(client, c.API))
	fmt.Printf("XAPI: %s [%s]\n", c.XAPI, pingUrl(client, c.XAPI))
	fmt.Printf("Exercises Directory: %s\n", c.Dir)
}

func checkLatestRelease(client http.Client) (*release, error) {
	resp, err := client.Get("https://api.github.com/repos/exercism/cli/releases/latest")
	if err != nil {
		return nil, errors.New("unable to get latest CLI release: " + err.Error())
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, errors.New("error decoding latest release response: " + err.Error())
	}

	return &rel, nil
}

func pingUrl(client http.Client, url string) string {
	res, err := client.Get(url)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	return "connected"
}
