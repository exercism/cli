package api

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const dateFormat = "January 2, 2006"

// StatusInfo contains information about a user's status on a particular language track.
type StatusInfo struct {
	TrackID         string `json:"track_id"`
	Recent          *Recent
	FetchedProblems *Slugs `json:"fetched"`
	SkippedProblems *Slugs `json:"skipped"`
}

// Recent contains information about the user's most recently submitted exercise on a particular language track.
type Recent struct {
	Problem     string `json:"problem"`
	SubmittedAt string `json:"submitted_at"`
}

// Slugs is a collection of slugs, all of which are the names of exercises.
type Slugs []string

func (r *Recent) String() string {
	submittedAt, err := time.Parse(time.RFC3339Nano, r.SubmittedAt)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(" - %s (submitted on %s)", r.Problem, submittedAt.Format(dateFormat))
}

func (s *StatusInfo) String() string {
	if len(*s.FetchedProblems) == 0 && s.Recent.Problem == "" {
		return fmt.Sprintf("\nYou have yet to begin the %s track!\n", s.TrackID)
	}

	msg := `
Your status on the %s track:

Most recently submitted exercise:
%s

Exercises fetched but not submitted:
%s

Exercises skipped:
%s
`

	return fmt.Sprintf(msg, s.TrackID, s.Recent, s.FetchedProblems, s.SkippedProblems)
}

func (s Slugs) String() string {
	for i, problem := range s {
		s[i] = fmt.Sprintf(" - %s", problem)
	}
	return strings.Join(s, "\n")
}
