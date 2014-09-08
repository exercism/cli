package handlers

import (
	"fmt"

	"github.com/exercism/cli/api"
)

type Item struct {
	*api.Problem
	dir string
}

func (it *Item) Path() string {
	return fmt.Sprintf("%s/%s", it.dir, it.Problem.ID)
}

func (it *Item) Save() error {
	return it.Problem.Save(it.dir)
}
