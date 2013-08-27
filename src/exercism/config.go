package exercism

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const FILENAME = ".exercism.go"

type Config struct {
	GithubUsername    string `json:"githubUsername"`
	ApiKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
	Hostname          string `json:"hostname"`
}

func ConfigFromFile(dir string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(configFilename(dir))
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}

func ConfigToFile(dir string, c Config) (err error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return
	}

	filename := configFilename(dir)
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return
	}
	fmt.Printf("Your credentials have been written to %s\n", filename)
	return
}

func configFilename(dir string) string {
	return dir + "/" + FILENAME
}
