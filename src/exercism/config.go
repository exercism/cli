package exercism

import (
	"encoding/json"
	"io/ioutil"
)

const FILENAME = ".exercism"

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

	err = ioutil.WriteFile(configFilename(dir), bytes, 0644)
	if err != nil {
		return
	}

	return
}

type Config struct {
	GithubUsername    string `json:"githubUsername"`
	ApiKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
}

func configFilename(dir string) string {
	return dir + "/" + FILENAME
}
