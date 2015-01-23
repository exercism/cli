package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Download returns specified submissions and related problem.
func Download(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	if len(args) != 1 {
		msg := "Usage: exercism download SUBMISSION_ID"
		log.Fatal(msg)
	}

	submission, err := client.Download(args[0])
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(c.Dir, "solutions", submission.Username, submission.Language, submission.Slug, args[0])

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}

	for name, contents := range submission.ProblemFiles {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, name), []byte(contents), 0755); err != nil {
			log.Fatal(err)
		}
	}

	for name, contents := range submission.SolutionFiles {
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, name), []byte(contents), 0755); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Successfully downloaded submission.\n\nThe submission can be viewed at:\n %s\n\n", path)

}
