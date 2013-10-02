package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func ConfigToFile(dir string, c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	filename := configFilename(dir)
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Your credentials have been written to %s\n", filename)
	return nil
}

func DemoDirectory() (string, error) {
	dir, err := os.Getwd()
	return fmt.Sprintf("%s/exercism-demo", dir), err
}

func configFilename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, FILENAME)
}
