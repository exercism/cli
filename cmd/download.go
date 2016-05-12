package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Download returns specified iteration with its related problem.
func Download(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: exercism download SUBMISSION_ID")
		os.Exit(1)
	}

	submission, err := client.Download(args[0])
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(c.Dir, "solutions", submission.Username, submission.TrackID, submission.Slug, args[0])

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}

	for name, contents := range submission.ProblemFiles {
		if err := writeFile(fmt.Sprintf("%s/%s", path, name), contents); err != nil {
			log.Fatalf("Unable to write file %s: %s", name, err)
		}
	}

	for name, contents := range submission.SolutionFiles {
		filename := strings.TrimPrefix(name, strings.ToLower("/"+submission.TrackID+"/"+submission.Slug+"/"))
		if err := writeFile(fmt.Sprintf("%s/%s", path, filename), contents); err != nil {
			log.Fatalf("Unable to write file %s: %s", name, err)
		}
	}

	fmt.Printf("Successfully downloaded submission.\n\nThe submission can be viewed at:\n %s\n\n", path)

	return nil

}

// writeFile writes the given contents to the given path, creating any necessary parent directories.
// This is useful because both problem files and solution files may have directory structures.
func writeFile(path, contents string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(contents), 0644)
}
