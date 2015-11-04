package cmd

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
)

// Fetch downloads exercism problems and writes them to disk.
func Fetch(ctx *cli.Context) {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	problems, err := client.Fetch(ctx.Args())
	if err != nil {
		log.Fatal(err)
	}

	submissionInfo, err := client.Submissions()
	if err != nil {
		log.Fatal(err)
	}

	if err := setSubmissionState(problems, submissionInfo); err != nil {
		log.Fatal(err)
	}

	hw := user.NewHomework(problems, c)
	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}

	hw.Summarize(user.HWAll)
}

func setSubmissionState(problems []*api.Problem, submissionInfo map[string][]api.SubmissionInfo) error {
	for _, problem := range problems {
		langSubmissions := submissionInfo[problem.TrackID]
		for _, submission := range langSubmissions {
			if submission.Slug == problem.Slug {
				problem.Submitted = true
			}
		}
	}

	return nil
}
