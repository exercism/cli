package user

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/exercism/cli/api"
)

// Item is a problem that has been fetched from the APIs.
// It contains some data specific to this particular request and user
// in order to give a useful report to the user about what has been fetched.
type Item struct {
	*api.Problem
	dir       string
	isNew     bool
	isUpdated bool
}

// Path is the location of this item on the user's filesystem.
func (it *Item) Path() string {
	return filepath.Join(it.dir, it.TrackID, it.Slug)
}

// Matches determines whether or not this item matches the given filter.
func (it *Item) Matches(filter HWFilter) bool {
	switch filter {
	case HWNew:
		return it.isNew
	case HWUpdated:
		return it.isUpdated
	}
	return true
}

// Save writes the embedded problem to the filesystem.
func (it *Item) Save() error {
	if _, err := os.Stat(it.Path()); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		it.isNew = true
	}

	for name, text := range it.Files {
		file := filepath.Join(it.Path(), name)

		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			return err
		}

		if _, err := os.Stat(file); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if !it.isNew {
				it.isUpdated = true
			}

			if err = ioutil.WriteFile(file, []byte(text), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
