package user

import (
	"fmt"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// HWFilter is used to categorize homework items.
type HWFilter int

// SummaryOption allows selective display of summary items.
type SummaryOption HWFilter

const (
	// HWAll represents all items in the collection.
	HWAll = iota
	// HWUpdated represents problems where files have been added.
	HWUpdated
	// HWNew represents newly fetched problems.
	HWNew
	// HWNotSubmitted represents problems that have not yet been submitted for review.
	HWNotSubmitted
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

	hw.template = "%s%s %s\n"
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
