package cmd

import (
	"log"
	"path/filepath"

	"github.com/urfave/cli"
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

	args := ctx.Args()
	problems, err := client.Fetch(args)

	if ctx.Bool("all") {
		if len(args) > 0 {
			trackID := args[0]
			problems = fetchAll(trackID, client)
		} else {
			log.Fatalf("You must supply a track to fetch all exercises")
		}
	}

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

func fetchAll(trackID string, client *api.Client) []*api.Problem {
	list, err := client.List(trackID)
	if err != nil {
		log.Fatal(err)
	}

	problems := make([]*api.Problem, len(list))
	for i, prob := range list {
		p, err := client.Fetch([]string{trackID, prob})
		if err != nil {
			log.Fatal(err)
		}
		problems[i] = p[0]
	}
	return problems
}
