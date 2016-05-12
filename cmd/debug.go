package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/paths"
)

type pingResult struct {
	URL     string
	Service string
	Status  string
	Latency time.Duration
}

// Debug provides information about the user's environment and configuration.
func Debug(ctx *cli.Context) error {
	defer fmt.Printf("\nIf you are having trouble and need to file a GitHub issue (https://github.com/exercism/exercism.io/issues) please include this information (except your API key. Keep that private).\n")

	client := &http.Client{Timeout: 20 * time.Second}

	fmt.Printf("\n**** Debug Information ****\n")
	fmt.Printf("Exercism CLI Version: %s\n", ctx.App.Version)

	rel, err := fetchLatestRelease(*client)
	if err != nil {
		log.Println("unable to fetch latest release: " + err.Error())
	} else {
		if rel.Version() != ctx.App.Version {
			defer fmt.Printf("\nA newer version of the CLI (%s) can be downloaded here: %s\n", rel.TagName, rel.Location)
		}
		fmt.Printf("Exercism CLI Latest Release: %s\n", rel.Version())
	}

	fmt.Printf("OS/Architecture: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Build OS/Architecture %s/%s\n", BuildOS, BuildARCH)
	if BuildARM != "" {
		fmt.Printf("Build ARMv%s\n", BuildARM)
	}

	fmt.Printf("Home Dir: %s\n", paths.Home)

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
		if c.APIKey != "" {
			fmt.Printf("API Key: %s\n", c.APIKey)
		} else {
			fmt.Println("API Key: Please set your API Key to access all of the CLI features")
		}
	} else {
		fmt.Println("Config file: <not configured>")
		fmt.Println("API Key: Please set your API Key to access all of the CLI features")
	}
	fmt.Printf("Exercises Directory: %s\n", c.Dir)

	fmt.Println("Testing API endpoints reachability")

	endpoints := map[string]string{
		"API":        c.API,
		"XAPI":       c.XAPI,
		"GitHub API": "https://api.github.com/",
	}

	var wg sync.WaitGroup
	results := make(chan pingResult)
	defer close(results)

	wg.Add(len(endpoints))

	for service, url := range endpoints {
		go func(service, url string) {
			now := time.Now()
			res, err := client.Get(url)
			delta := time.Since(now)
			if err != nil {
				results <- pingResult{
					URL:     url,
					Service: service,
					Status:  err.Error(),
					Latency: delta,
				}
				return
			}
			defer res.Body.Close()

			results <- pingResult{
				URL:     url,
				Service: service,
				Status:  "connected",
				Latency: delta,
			}
		}(service, url)
	}

	go func() {
		for r := range results {
			fmt.Printf(
				"\t* %s: %s [%s] %s\n",
				r.Service,
				r.URL,
				r.Status,
				r.Latency,
			)
			wg.Done()
		}
	}()

	wg.Wait()

	return nil
}
