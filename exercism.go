package main

import (
	"fmt"
	"github.com/exercism/cli/configuration"
	"os"
	"path/filepath"
)

func logout(path string) {
	os.Remove(path)
}

func absolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

func askForConfigInfo() (c configuration.Config) {
	var un, key, dir string

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Print("Your GitHub username: ")
	_, err = fmt.Scanln(&un)
	if err != nil {
		panic(err)
	}

	fmt.Print("Your exercism.io API key: ")
	_, err = fmt.Scanln(&key)
	if err != nil {
		panic(err)
	}

	fmt.Println("What is your exercism exercises project path?")
	fmt.Printf("Press Enter to select the default (%s):\n", currentDir)
	fmt.Print("> ")
	_, err = fmt.Scanln(&dir)
	if err != nil && err.Error() != "unexpected newline" {
		panic(err)
	}

	if dir == "" {
		dir = currentDir
	}

	dir = configuration.ReplaceTilde(dir)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = fmt.Errorf("Error making directory %v: [%v]", dir, err)
		return
	}

	dir, err = absolutePath(dir)
	if err != nil {
		panic(err)
	}

	return configuration.Config{GithubUsername: un, ApiKey: key, ExercismDirectory: dir, Hostname: "http://exercism.io"}
}
