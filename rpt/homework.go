package rpt

import (
	"fmt"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// HWFilter is used to categorize homework items.
type HWFilter int

const (
	// HWAll represents all items in the collection.
	HWAll = iota
	// HWUpdated represents problems that were already on the
	// user's filesystem, where one or more new files have been added.
	HWUpdated
	// HWNew represents problems that did not yet exist on the
	// user's filesystem.
	HWNew
)

// Homework is a collection of problems that were fetched from the APIs.
type Homework struct {
	Items    []*Item
	template string
}

// NewHomework decorates a problem set with some additional data based on the
// user's system.
func NewHomework(problems []*api.Problem, c *config.Config) *Homework {
	hw := Homework{}
	for _, problem := range problems {
		item := &Item{
			Problem: problem,
			dir:     c.Dir,
		}
		hw.Items = append(hw.Items, item)
	}

	hw.template = fmt.Sprintf("%%%ds %%s\n", hw.maxTitleWidth())
	return &hw
}

// Save saves all problems in the problem set.
func (hw *Homework) Save() error {
	for _, item := range hw.Items {
		if err := item.Save(); err != nil {
			return err
		}
	}
	return nil
}

// ItemsMatching returns a subset of the set of problems.
func (hw *Homework) ItemsMatching(filter HWFilter) []*Item {
	items := []*Item{}
	for _, item := range hw.Items {
		if item.Matches(filter) {
			items = append(items, item)
		}
	}
	return items
}

// Report outputs a list of the problems in the set.
// It prints the track name, the problem name, and the full
// path to the problem on the user's filesystem.
func (hw *Homework) Report(filter HWFilter) {
	items := hw.ItemsMatching(filter)
	hw.heading(filter, len(items))
	for _, item := range items {
		fmt.Printf(hw.template, item.String(), item.Path())
	}
}

func (hw *Homework) heading(filter HWFilter, count int) {
	if count == 0 {
		return
	}
	fmt.Println()

	if filter == HWAll {
		return
	}

	unit := "problems"
	if count == 1 {
		unit = "problem"
	}

	var status string
	switch filter {
	case HWUpdated:
		status = "Updated:"
	case HWNew:
		status = "New:"
	}
	summary := fmt.Sprintf("%d %s", count, unit)
	fmt.Printf(hw.template, status, summary)
}

func (hw *Homework) maxTitleWidth() int {
	var max int
	for _, item := range hw.Items {
		if len(item.String()) > max {
			max = len(item.String())
		}
	}
	return max
}

// Summarize prints a full report of new and updated items in the set.
func (hw *Homework) Summarize() {
	hw.Report(HWUpdated)
	hw.Report(HWNew)

	fresh := len(hw.ItemsMatching(HWNew))
	updated := len(hw.ItemsMatching(HWUpdated))
	unchanged := len(hw.Items) - updated - fresh
	fmt.Printf("\nunchanged: %d, updated: %d, new: %d\n\n", unchanged, updated, fresh)
}
