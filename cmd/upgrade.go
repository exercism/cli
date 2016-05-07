package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/kardianos/osext"
)

var (
	// BuildOS is the operating system (GOOS) used during the build process.
	BuildOS string
	// BuildARM is the ARM version (GOARM) used during the build process.
	BuildARM string
	// BuildARCH is the architecture (GOARCH) used during the build process.
	BuildARCH string
)

// Upgrade allows the user to upgrade to the latest version of the CLI.
func Upgrade(ctx *cli.Context) error {
	client := http.Client{Timeout: 10 * time.Second}
	rel, err := fetchLatestRelease(client)
	if err != nil {
		log.Fatal("unable to check latest release version: " + err.Error())
	}

	// TODO: This checks strings against string
	// Version 2.2.0 against release 2.1.0
	// will trigger an upgrade...
	// should probably parse semver and compare
	if rel.Version() == ctx.App.Version {
		fmt.Println("Your CLI is up to date!")
		return nil
	}

	// Locate the path to the current `exercism` executable.
	dest, err := osext.Executable()
	if err != nil {
		log.Fatalf("Unable to find current executable path: %s", err)
	}

	var (
		OS   = osMap[runtime.GOOS]
		ARCH = archMap[runtime.GOARCH]
	)

	if OS == "" || ARCH == "" {
		log.Fatalf("unable to upgrade: OS %s ARCH %s", OS, ARCH)
	}

	buildName := fmt.Sprintf("%s-%s", OS, ARCH)
	if BuildARCH == "arm" {
		if BuildARM == "" {
			log.Fatalf("unable to upgrade: arm version not found")
		}
		buildName = fmt.Sprintf("%s-v%s", buildName, BuildARM)
	}

	var downloadRC *bytes.Reader
	for _, a := range rel.Assets {
		if strings.Contains(a.Name, buildName) {
			fmt.Printf("Downloading %s\n", a.Name)
			downloadRC, err = a.download()
			if err != nil {
				log.Fatalf("error downloading executable: %s\n", err)
			}
			break
		}
	}
	if downloadRC == nil {
		log.Fatalf("No executable found for %s/%s%s", BuildOS, BuildARCH, BuildARM)
	}

	if OS == "windows" {
		err = installZip(downloadRC, dest)
	} else {
		err = installTgz(downloadRC, dest)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully upgraded!")

	return nil
}
