package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
	"github.com/robphoenix/cli/user"
	"github.com/urfave/cli"
)

// Fetch downloads exercism exercises and writes them to disk.
func Fetch(ctx *cli.Context) error {
	c, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(c)

	args := ctx.Args()
	var exercises []*api.Exercise

	if ctx.Bool("all") {
		if len(args) > 0 {
			trackID := args[0]
			fmt.Printf("\nFetching all exercises for the %s track...\n\n", trackID)
			p, err := client.FetchAll(trackID)
			if err != nil {
				log.Fatal(err)
			}
			exercises = p
		} else {
			log.Fatalf("You must supply a track to fetch all exercises")
		}
	} else {
		p, err := client.Fetch(args)
		if err != nil {
			log.Fatal(err)
		}
		exercises = p
	}

	submissionInfo, err := client.Submissions()
	if err != nil {
		log.Fatal(err)
	}

	if err := setSubmissionState(exercises, submissionInfo); err != nil {
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
	hw := user.NewHomework(exercises, c)

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

func setSubmissionState(exercises []*api.Exercise, submissionInfo map[string][]api.SubmissionInfo) error {
	for _, exercise := range exercises {
		langSubmissions := submissionInfo[exercise.TrackID]
		for _, submission := range langSubmissions {
			if submission.Slug == exercise.Slug {
				exercise.Submitted = true
			}
		}
	}

	return nil
}
