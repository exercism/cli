package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/kardianos/osext"
)

var (
	// BuildOS is the GOOS used during the build process
	BuildOS string
	// Build ARMV is the GOARM used during the build process
	BuildARM string
	// BuildARCH is the GOARCH used during the build process
	BuildARCH string
)

// Upgrade command allows the user to upgrade to the latest CLI version
func Upgrade(ctx *cli.Context) {
	client := http.Client{Timeout: 5 * time.Second}
	rel, err := checkLatestRelease(client)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: This checks strings against string
	// Version 2.2.0 against release 2.1.0
	// will trigger an upgrade...
	// should probably parse semver and compare
	if rel.Version() == ctx.App.Version {
		fmt.Println("Your CLI is up to date!")
		return
	}

	// current executable path
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
			downloadRC, err = a.Download()
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
}

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

const installFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC

func installTgz(source *bytes.Reader, dest string) error {
	gr, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileCopy, err := os.OpenFile(dest, installFlag, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		if _, err = io.Copy(fileCopy, tr); err != nil {
			return err
		}
	}

	return nil
}

func installZip(source *bytes.Reader, dest string) error {
	zr, err := zip.NewReader(source, int64(source.Len()))
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		fileCopy, err := os.OpenFile(dest, installFlag, f.Mode())
		if err != nil {
			return err
		}
		defer fileCopy.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(fileCopy, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

type asset struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
}

func (a *asset) Download() (*bytes.Reader, error) {
	downloadURL := fmt.Sprintf("https://api.github.com/repos/exercism/cli/releases/assets/%d", a.ID)
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, err
	}
	// https://developer.github.com/v3/repos/releases/#get-a-single-release-asset
	req.Header.Set("Accept", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bs), nil
}

type release struct {
	Location string  `json:"html_url"`
	TagName  string  `json:"tag_name"`
	Assets   []asset `json:"assets"`
}

func (r *release) Version() string {
	return strings.TrimPrefix(r.TagName, "v")
}
