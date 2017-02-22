package user

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
)

// HWFilter is used to categorize homework items.
type HWFilter int

// SummaryOption allows selective display of summary items.
type SummaryOption HWFilter

const (
	// HWAll represents all items in the collection.
	HWAll = iota
	// HWUpdated represents exercises where files have been added.
	HWUpdated
	// HWNew represents newly fetched exercises.
	HWNew
	// HWNotSubmitted represents exercises that have not yet been submitted for review.
	HWNotSubmitted
)

// Homework is a collection of exercises that were fetched from the APIs.
type Homework struct {
	Items    []*Item
	template string
}

// NewHomework decorates an exercise set with some additional data based on the
// user's system.
func NewHomework(exercises []*api.Exercise, c *config.Config) *Homework {
	hw := Homework{}
	for _, exercise := range exercises {
		item := &Item{
			Exercise: exercise,
			dir:      c.Dir,
		}
		hw.Items = append(hw.Items, item)
	}

	hw.template = "%s%s %s\n"
	return &hw
}

// Save saves all exercises in the exercise set.
func (hw *Homework) Save() error {
	for _, item := range hw.Items {
		if err := item.Save(); err != nil {
			return err
		}
	}
	return nil
}

// RejectMissingTracks removes any items that are part of tracks the user
// doesn't currently have a folder for on their local machine. This
// only happens when a user calls `exercism fetch` without any arguments.
func (hw *Homework) RejectMissingTracks(dirMap map[string]bool) error {
	items := []*Item{}
	for _, item := range hw.Items {
		dir := filepath.Join(item.dir, item.TrackID)
		if dirMap[dir] {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return fmt.Errorf(`
You have yet to start a language track!
View all available language tracks with "exercism tracks"
Fetch exercises for your first track with "exercism fetch TRACK_ID"`)
	}
	hw.Items = items
	return nil
}

// ItemsMatching returns a subset of the set of exercises.
func (hw *Homework) ItemsMatching(filter HWFilter) []*Item {
	items := []*Item{}
	for _, item := range hw.Items {
		if item.Matches(filter) {
			items = append(items, item)
		}
	}
	return items
}

// Report outputs a list of the exercises in the set.
// It prints the track name, the exercise name, and the full
// path to the exercise on the user's filesystem.
func (hw *Homework) Report(filter HWFilter) {
	if hw == nil {
		return
	}
	width := hw.maxTitleWidth()
	items := hw.ItemsMatching(filter)
	hw.heading(filter, len(items), width)
	for _, item := range items {
		fmt.Print(item.Report(hw.template, width))
	}
}

func (hw *Homework) heading(filter HWFilter, count, width int) {
	if count == 0 {
		return
	}
	fmt.Println()

	if filter == HWAll {
		return
	}

	unit := "exercises"
	if count == 1 {
		unit = "exercise"
	}

	var status string
	switch filter {
	case HWUpdated:
		status = "Updated:"
	case HWNew:
		status = "New:"
	case HWNotSubmitted:
		status = "Not Submitted:"
	}
	summary := fmt.Sprintf("%d %s", count, unit)
	var padding string
	if width > len(status) {
		padding = strings.Repeat(" ", width-len(status))
	}
	fmt.Printf(hw.template, status, padding, summary)
}

func (hw *Homework) maxTitleWidth() int {
	if hw == nil {
		return 0
	}
	var width int
	for _, item := range hw.Items {
		if len(item.String()) > width {
			width = len(item.String())
		}
	}
	return width
}

// Summarize prints a full report of new and updated items in the set.
func (hw *Homework) Summarize(summaryFilter SummaryOption) {
	hw.Report(HWUpdated)

	if summaryFilter != HWNotSubmitted {
		hw.Report(HWNotSubmitted)
	}

	hw.Report(HWNew)

	fresh := len(hw.ItemsMatching(HWNew))
	updated := len(hw.ItemsMatching(HWUpdated))
	unchanged := len(hw.Items) - updated - fresh
	fmt.Printf("\nunchanged: %d, updated: %d, new: %d\n\n", unchanged, updated, fresh)
}
