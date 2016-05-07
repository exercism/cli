package cmd

import (
	"log"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/user"
)

// Fetch downloads exercism problems and writes them to disk.
func Fetch(ctx *cli.Context) error {
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

	dirs, err := filepath.Glob(filepath.Join(c.Dir, "*"))
	if err != nil {
		log.Fatal(err)
	}

	dirMap := make(map[string]bool)
	for _, dir := range dirs {
		dirMap[dir] = true
	}
	hw := user.NewHomework(problems, c)

	if len(ctx.Args()) == 0 {
		if err := hw.RejectMissingTracks(dirMap); err != nil {
			log.Fatal(err)
		}
	}

	if err := hw.Save(); err != nil {
		log.Fatal(err)
	}

	hw.Summarize(user.HWAll)

	return nil
	// return cli.NewExitError("no good", 10)
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
