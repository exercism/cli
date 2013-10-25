package configuration

import (
	"fmt"
	"os"
)

func Demo() (c Config, err error) {
	demoDir, err := demoDirectory()
	if err != nil {
		return
	}
	c = Config{
		Hostname:          "http://exercism.io",
		ApiKey:            "",
		ExercismDirectory: demoDir,
	}
	return
}

func demoDirectory() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return
	}
	dir = fmt.Sprintf("%s/exercism-demo", dir)
	return
}
