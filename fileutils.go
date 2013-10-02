package main

import (
	"os"
	"runtime"
	"strings"
)

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
