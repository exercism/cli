package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const FILENAME = ".exercism.go"

type Config struct {
	GithubUsername    string `json:"githubUsername"`
	ApiKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
	Hostname          string `json:"hostname"`
}

func ToFile(path string, c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Your credentials have been written to %s\n", path)
	return nil
}

func ReplaceTilde(oldPath string) string {
	return strings.Replace(oldPath, "~/", HomeDir()+"/", 1)
}

// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}
func FromFile(path string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}

func Filename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, FILENAME)
}
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
	dir = filepath.Join(dir, "exercism-demo")
	return
}
