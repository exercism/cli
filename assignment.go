package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Assignment struct {
	Track   string
	Slug    string
	Files   map[string]string
	IsFresh bool `json:"fresh"`
}

func SaveAssignment(dir string, a Assignment) (err error) {
	root := fmt.Sprintf("%s/%s/%s", dir, a.Track, a.Slug)

	for name, text := range a.Files {
		file := fmt.Sprintf("%s/%s", root, name)
		dir := filepath.Dir(file)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			err = fmt.Errorf("Error making directory %v: [%v]", dir, err)
			return
		}
		if _, err = os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				err = ioutil.WriteFile(file, []byte(text), 0644)
				if err != nil {
					err = fmt.Errorf("Error writing file %v: [%v]", name, err)
					return
				}
			}
		}
	}

	fresh := " "
	if a.IsFresh {
		fresh = "*"
	}
	fmt.Println(fresh, a.Track, "-", a.Slug)

	return
}
