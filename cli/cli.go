package cli

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/exercism/cli/debug"
	update "github.com/inconshreveable/go-update"
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

var (
	// HTTPClient is the client used to make HTTP calls in the cli package.
	HTTPClient = &http.Client{Timeout: 10 * time.Second}
	// LatestReleaseURL is the endpoint that provides information about the latest release.
	LatestReleaseURL = "https://api.github.com/repos/exercism/cli/releases/latest"
)

// CLI is information about the CLI itself.
type CLI struct {
	Version       string
	LatestRelease *Release
}

// New creates a CLI, setting it to a particular version.
func New(version string) *CLI {
	return &CLI{
		Version: version,
	}
}

// IsUpgradeNeeded compares the current version to that of the latest release.
func (c *CLI) IsUpgradeNeeded() (bool, error) {
	if c.LatestRelease == nil {
		resp, err := HTTPClient.Get(LatestReleaseURL)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		var rel Release
		if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
			return false, err
		}
		c.LatestRelease = &rel
	}

	rv, err := semver.Make(c.LatestRelease.Version())
	if err != nil {
		return false, fmt.Errorf("unable to parse latest version (%s): %s", c.LatestRelease.Version(), err)
	}
	cv, err := semver.Make(c.Version)
	if err != nil {
		return false, fmt.Errorf("unable to parse current version (%s): %s", c.Version, err)
	}

	return rv.GT(cv), nil
}

// Upgrade allows the user to upgrade to the latest version of the CLI.
func (c *CLI) Upgrade() error {
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
	for _, a := range c.LatestRelease.Assets {
		if strings.Contains(a.Name, buildName) {
			debug.Printf("Downloading %s\n", a.Name)
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

	bin, err := extractBinary(downloadRC, OS)
	if err != nil {
		return err
	}
	defer bin.Close()

	return update.Apply(bin, update.Options{})
}

func extractBinary(source *bytes.Reader, os string) (binary io.ReadCloser, err error) {
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
			}
			if _, err := tmpfile.Seek(0, 0); err != nil {
				return nil, err
			}

			binary = tmpfile
		}
	}

	return binary, nil
}
