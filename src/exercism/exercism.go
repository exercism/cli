package exercism

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const FILENAME = ".exercism"

func Login(configDir string, config Config) (err error) {
	contents, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(createFilename(configDir), contents, 0644)
	return err
}

func Logout(configDir string) {
	os.Remove(createFilename(configDir))
}

func createFilename(dir string) string {
	return dir + "/" + FILENAME
}
