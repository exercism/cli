package user

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/robphoenix/cli/api"
)

// Item is an exercise that has been fetched from the APIs.
// It contains some data specific to this particular request and user
// in order to give a useful report to the user about what has been fetched.
type Item struct {
	*api.Exercise
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
	case HWNotSubmitted:
		return !it.Submitted
	}
	return true
}

// Save writes the embedded exercise to the filesystem.
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

			if runtime.GOOS == "windows" {
				text = strings.Replace(text, "\n", "\r\n", -1)
			}

			if err := ioutil.WriteFile(file, []byte(text), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// Report outputs the line's string and path in the format of the passed in template.
func (it *Item) Report(template string, max int) string {
	padding := strings.Repeat(" ", max-len(it.String()))
	return fmt.Sprintf(template, it.String(), padding, it.Path())
}
