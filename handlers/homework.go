package handlers

import (
	"fmt"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

type Item struct {
	Title string
	Path  string
}

type Homework struct {
	Items    []*Item
	template string
}

func (hw Homework) MaxTitleWidth() int {
	var max int
	for _, item := range hw.Items {
		if len(item.Title) > max {
			max = len(item.Title)
		}
	}
	return max
}

func NewHomework(problems []*api.Problem, c *config.Config) Homework {
	hw := Homework{}
	for _, problem := range problems {
		item := &Item{
			Title: problem.String(),
			Path:  fmt.Sprintf("%s/%s", c.Dir, problem.ID),
		}
		hw.Items = append(hw.Items, item)
	}

	hw.template = fmt.Sprintf("%%%ds %%s\n", hw.MaxTitleWidth())
	return hw
}

func (hw Homework) Report() {
	for _, item := range hw.Items {
		fmt.Printf(hw.template, item.Title, item.Path)
	}
}
