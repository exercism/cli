package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
	for name, text := range it.Files {
		file := fmt.Sprintf("%s/%s/%s", it.dir, it.ID, name)

		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			return err
		}

		if _, err := os.Stat(file); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			err = ioutil.WriteFile(file, []byte(text), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
