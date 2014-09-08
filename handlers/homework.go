package handlers

import (
	"fmt"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

type HWFilter int

const (
	HWAll = iota
	HWUpdated
	HWNew
)

type Homework struct {
	Items    []*Item
	template string
}

func NewHomework(problems []*api.Problem, c *config.Config) *Homework {
	hw := Homework{}
	for _, problem := range problems {
		item := &Item{
			Problem: problem,
			dir:     c.Dir,
		}
		hw.Items = append(hw.Items, item)
	}

	hw.template = fmt.Sprintf("%%%ds %%s\n", hw.MaxTitleWidth())
	return &hw
}

func (hw *Homework) Save() error {
	for _, item := range hw.Items {
		err := item.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (hw *Homework) ItemsMatching(filter HWFilter) []*Item {
	items := []*Item{}
	for _, item := range hw.Items {
		if item.Matches(filter) {
			items = append(items, item)
		}
	}
	return items
}

func (hw *Homework) Report(filter HWFilter) {
	items := hw.ItemsMatching(filter)
	hw.Heading(filter, len(items))
	for _, item := range items {
		fmt.Printf(hw.template, item.String(), item.Path())
	}
}

func (hw *Homework) Heading(filter HWFilter, count int) {
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

func (hw *Homework) MaxTitleWidth() int {
	var max int
	for _, item := range hw.Items {
		if len(item.String()) > max {
			max = len(item.String())
		}
	}
	return max
}

func (hw *Homework) Summarize() {
	hw.Report(HWUpdated)
	hw.Report(HWNew)

	fresh := len(hw.ItemsMatching(HWNew))
	updated := len(hw.ItemsMatching(HWUpdated))
	unchanged := len(hw.Items) - updated - fresh
	fmt.Printf("\nunchanged: %d, updated: %d, new: %d\n\n", unchanged, updated, fresh)
}
