package cli

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		"darwin":  "darwin",
		"freebsd": "freebsd",
		"linux":   "linux",
		"openbsd": "openbsd",
		"windows": "windows",
	}

	archMap = map[string]string{
		"386":   "i386",
		"amd64": "x86_64",
		"arm":   "arm",
		"ppc64": "ppc64",
	}
)

var (
	// TimeoutInSeconds is the timeout the default HTTP client will use.
	TimeoutInSeconds = 60
	// HTTPClient is the client used to make HTTP calls in the cli package.
	HTTPClient = &http.Client{Timeout: time.Duration(TimeoutInSeconds) * time.Second}
	// ReleaseURL is the endpoint that provides information about cli releases.
	ReleaseURL = "https://api.github.com/repos/exercism/cli/releases"
)

// Updater is a simple upgradable file interface.
type Updater interface {
	IsUpToDate() (bool, error)
	Upgrade() error
}

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

// IsUpToDate compares the current version to that of the latest release.
func (c *CLI) IsUpToDate() (bool, error) {
	if c.LatestRelease == nil {
		if err := c.fetchLatestRelease(); err != nil {
			return false, err
		}
	}

	rv, err := semver.Make(c.LatestRelease.Version())
	if err != nil {
		return false, fmt.Errorf("unable to parse latest version (%s): %s", c.LatestRelease.Version(), err)
	}
	cv, err := semver.Make(c.Version)
	if err != nil {
		return false, fmt.Errorf("unable to parse current version (%s): %s", c.Version, err)
	}

	return cv.GTE(rv), nil
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

func (c *CLI) fetchLatestRelease() error {
	latestReleaseURL := fmt.Sprintf("%s/%s", ReleaseURL, "latest")
	resp, err := HTTPClient.Get(latestReleaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		msg := "failed to get the latest release\n"
		for k, v := range resp.Header {
			msg += fmt.Sprintf("\n  %s:\n    %s", k, v)
		}
		return fmt.Errorf(msg)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return err
	}
	c.LatestRelease = &rel
	return nil
}

func extractBinary(source *bytes.Reader, platform string) (binary io.ReadCloser, err error) {
	if platform == "windows" {
		zr, err := zip.NewReader(source, int64(source.Len()))
		if err != nil {
			return nil, err
		}

		for _, f := range zr.File {
			info := f.FileInfo()
			if info.IsDir() || !strings.HasSuffix(f.Name, ".exe") {
				continue
			}
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
			tmpfile, err := os.CreateTemp("", "temp-exercism")
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
