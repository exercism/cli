package exercism

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"strings"
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

func ConfigToFile(u user.User, dir string, c Config) (err error) {
	expandedPath, err := replaceTilde(u, c.ExercismDirectory)
	if err != nil {
		return
	}

	c.ExercismDirectory = expandedPath
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
	Hostname          string `json:"hostname"`
}

func configFilename(dir string) string {
	return dir + "/" + FILENAME
}

func replaceTilde(u user.User, oldPath string) (newPath string, err error) {
	dir := u.HomeDir
	newPath = strings.Replace(oldPath, "~/", dir+"/", 1)
	return
}
