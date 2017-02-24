package user

import (
	"fmt"
	"strings"

	"github.com/robphoenix/cli/api"
)

// Status is the status of a track (active/inactive).
type Status bool

const (
	// TrackActive represents an active track.
	// Exercises from active tracks will be delivered with the `fetch` command.
	TrackActive Status = true
	// TrackInactive represents an inactive track.
	// It is possible to fetch exercises from an inactive track, and
	// submit them to the website, but these will not automatically be
	// delivered in the global `fetch` command.
	TrackInactive Status = false
)

// Curriculum is a collection of language tracks.
type Curriculum struct {
	Tracks []*api.Track
	wLang  int
	wID    int
}

// NewCurriculum returns a collection of language tracks.
func NewCurriculum(tracks []*api.Track) *Curriculum {
	return &Curriculum{Tracks: tracks}
}

// Report creates a table of the tracks that have the requested status.
func (cur *Curriculum) Report(status Status) {
	for _, track := range cur.Tracks {
		if Status(track.Active) == status {
			fmt.Println(
				"    ",
				track.Language,
				strings.Repeat(" ", cur.lenLang()-len(track.Language)+1),
				track.ID,
				strings.Repeat(" ", cur.lenID()-len(track.ID)+1),
				track.Len(),
				"exercises",
			)
		}
	}
}

func (cur *Curriculum) lenLang() int {
	if cur.wLang > 0 {
		return cur.wLang
	}

	for _, track := range cur.Tracks {
		if len(track.Language) > cur.wLang {
			cur.wLang = len(track.Language)
		}
	}
	return cur.wLang
}

func (cur *Curriculum) lenID() int {
	if cur.wID > 0 {
		return cur.wID
	}

	for _, track := range cur.Tracks {
		if len(track.ID) > cur.wID {
			cur.wID = len(track.ID)
		}
	}
	return cur.wID
}
