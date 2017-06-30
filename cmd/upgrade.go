package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
	"github.com/urfave/cli"
)

var (
	// BuildOS is the operating system (GOOS) used during the build process.
	BuildOS string
	// BuildARM is the ARM version (GOARM) used during the build process.
	BuildARM string
	// BuildARCH is the architecture (GOARCH) used during the build process.
	BuildARCH string
)

var (
	osMap = map[string]string{
		"darwin":  "mac",
		"linux":   "linux",
		"windows": "windows",
	}

	archMap = map[string]string{
		"amd64": "64bit",
		"386":   "32bit",
		"arm":   "arm",
	}
)

type upgrader struct {
	client  *http.Client
	release *release
}

func (u *upgrader) fetchLatestRelease() (*release, error) {
	resp, err := u.client.Get("https://api.github.com/repos/exercism/cli/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	return &rel, nil
}

func (u *upgrader) IsUpgradeNeeded(currentVersion string) (bool, error) {
	rel := u.release
	latestVer, err := semver.Make(rel.Version())
	if err != nil {
		return false, fmt.Errorf("Unable to parse latest version (%s): %s", rel.Version(), err)
	}
	currentVer, err := semver.Make(currentVersion)
	if err != nil {
		return false, fmt.Errorf("Unable to parse current version (%s): %s", currentVersion, err)
	}

	return latestVer.GT(currentVer), nil
}

func NewUpgrader(client *http.Client) (*upgrader, error) {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	u := &upgrader{client: client}
	rel, err := u.fetchLatestRelease()
	if err != nil {
		return nil, err
	}
	u.release = rel
	return u, nil
}

func (u *upgrader) Upgrade() error {
	var (
		OS   = osMap[runtime.GOOS]
		ARCH = archMap[runtime.GOARCH]
	)

	if OS == "" || ARCH == "" {
		return fmt.Errorf("unable to upgrade: OS %s ARCH %s", OS, ARCH)
	}

	buildName := fmt.Sprintf("%s-%s", OS, ARCH)
	if BuildARCH == "arm" {
		if BuildARM == "" {
			return fmt.Errorf("unable to upgrade: arm version not found")
		}
		buildName = fmt.Sprintf("%s-v%s", buildName, BuildARM)
	}

	var downloadRC *bytes.Reader
	for _, a := range u.release.Assets {
		if strings.Contains(a.Name, buildName) {
			// TODO: This should be debug
			fmt.Printf("Downloading %s\n", a.Name)
			var err error
			downloadRC, err = a.download()
			if err != nil {
				return fmt.Errorf("error downloading executable: %s", err)
			}
			break
		}
	}
	if downloadRC == nil {
		return fmt.Errorf("no executable found for %s/%s%s", BuildOS, BuildARCH, BuildARM)
	}

	bin, err := u.extractBinary(downloadRC, OS)
	if err != nil {
		return err
	}
	defer bin.Close()

	if err := update.Apply(bin, update.Options{}); err != nil {
		return err
	}
	return nil
}

func (u *upgrader) extractBinary(source *bytes.Reader, os string) (binary io.ReadCloser, err error) {
	if os == "windows" {
		zr, err := zip.NewReader(source, int64(source.Len()))
		if err != nil {
			return nil, err
		}

		for _, f := range zr.File {
			return f.Open()
		}
	} else {
		gr, err := gzip.NewReader(source)
		if err != nil {
			return nil, err
		}
		defer gr.Close()

		tr := tar.NewReader(gr)
		for {
			_, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			tmpfile, err := ioutil.TempFile("", "temp-exercism")
			if err != nil {
				return nil, err
			}

			if _, err = io.Copy(tmpfile, tr); err != nil {
				return nil, err
			} else {
				if _, err := tmpfile.Seek(0, 0); err != nil {
					return nil, err
				}

				binary = tmpfile
			}
		}
	}

	return binary, nil
}

// Upgrade allows the user to upgrade to the latest version of the CLI.
func Upgrade(ctx *cli.Context) error {
	u, err := NewUpgrader(nil)
	if err != nil {
		log.Fatal(err)
	}

	upgradeNeeded, err := u.IsUpgradeNeeded(ctx.App.Version)
	if err != nil {
		log.Fatalf("unable to check for upgrade: %s", err)
		return err
	}

	if !upgradeNeeded {
		fmt.Println("Your CLI is up to date!")
		return nil
	}

	if err := u.Upgrade(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully upgraded!")
	return nil
}
