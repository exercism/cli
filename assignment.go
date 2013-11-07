package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Assignment struct {
	Track string
	Slug  string
	Files map[string]string
}

func SaveAssignment(dir string, a Assignment) (err error) {
	assignmentPath := fmt.Sprintf("%s/%s/%s", dir, a.Track, a.Slug)
	err = os.MkdirAll(assignmentPath, 0744)
	if err != nil {
		err = fmt.Errorf("Error creating assignment directory: [%s]", err)
		return
	}

	for name, text := range a.Files {
		filePath := fmt.Sprintf("%s/%s", assignmentPath, name)
		err = ioutil.WriteFile(filePath, []byte(text), 0644)
		if err != nil {
			err = fmt.Errorf("Error writing %v file: [%v]", name, err)
			return
		}
	}

	fmt.Println(a.Track, "-", a.Slug)

	return
}
