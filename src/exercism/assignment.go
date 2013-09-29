package exercism

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Assignment struct {
	Track    string
	Slug     string
	Readme   string
	TestFile string `json:"test_file"`
	Tests    string
}

func SaveAssignment(dir string, a Assignment) (err error) {
	assignmentPath := fmt.Sprintf("%s/%s/%s", dir, a.Track, a.Slug)
	err = os.MkdirAll(assignmentPath, 0744)
	if err != nil {
		err = fmt.Errorf("Error creating assignment directory: [%s]", err.Error())
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", assignmentPath, "README.md"), []byte(a.Readme), 0644)
	if err != nil {
		err = fmt.Errorf("Error writing README.md file: [%s]", err.Error())
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", assignmentPath, a.TestFile), []byte(a.Tests), 0644)
	if err != nil {
		err = fmt.Errorf("Error writing file %s: [%s]", a.TestFile, err.Error())
	}

	fmt.Println(a.Track, "-", a.Slug)

	return
}
