package cmd

import (
	"fmt"
	"log"
	"runtime"
	"os"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Gets a given user's particular submission.
func Download(ctx *cli.Context) { //Restore.go code.. for reference (we will modify)
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}

	args := ctx.Args()
	var url string
	switch len(args) {
	case 1:
	// Receives data from http://exercism.io/api/v1/submissions/:submissionKey
		url = fmt.Sprintf("%s/api/v1/submissions/%s", c.API, args[0])
	default:
		msg := "Usage: exercism download\n		or: exercism download SUBMISSION_ID"
		log.Fatal(msg)
	}

	submission, err := api.Download(url)

	if err != nil {
		log.Fatal(err)
	}

	var path string

	if runtime.GOOS == "windows" {
		path = fmt.Sprintf("%s\\solutions\\%s\\%s\\%s\\%s\\", c.Dir, submission.UserName, submission.Language, submission.Slug, args[0])
	} else {
		path = fmt.Sprintf("%s/solutions/%s/%s/%s/%s/", c.Dir, submission.UserName, submission.Language, submission.Slug, args[0])
	}

	// if err := os.RemoveAll(path); err != nil {
	// 	log.Fatal(err)
	// }

	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}

	for k := range submission.Problem {
		if err := ioutil.WriteFile(fmt.Sprintf("%s%s", path, k), []byte(submission.Problem[k]), 0755); err != nil {
			log.Fatal(err)
		}
	}

	for k := range submission.Code {
		if err := ioutil.WriteFile(fmt.Sprintf("%s%s", path, k), []byte(submission.Code[k]), 0755); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Successfully downloaded submission.\n\nThe submission can be viewed at:\n %s\n\n", path)

}
