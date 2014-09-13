package handlers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/exercism/cli/api"
)

type Item struct {
	*api.Problem
	dir       string
	isNew     bool
	isUpdated bool
}

func (it *Item) Path() string {
	return filepath.Join(it.dir, it.TrackID, it.Slug)
}

func (it *Item) Matches(filter HWFilter) bool {
	switch filter {
	case HWNew:
		return it.isNew
	case HWUpdated:
		return it.isUpdated
	}
	return true
}

func (it *Item) Save() error {
	if _, err := os.Stat(it.Path()); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		it.isNew = true
	}

	for name, text := range it.Files {
		file := filepath.Join(it.Path(), name)

		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return err
		}

		if _, err := os.Stat(file); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if !it.isNew {
				it.isUpdated = true
			}

			err = ioutil.WriteFile(file, []byte(text), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
