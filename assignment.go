package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Assignment struct {
	Track string
	Slug  string
	Files map[string]string
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
		err = ioutil.WriteFile(file, []byte(text), 0644)
		if err != nil {
			err = fmt.Errorf("Error writing file %v: [%v]", name, err)
			return
		}
	}

	fmt.Println(a.Track, "-", a.Slug)

	return
}
