package main

import (
	"github.com/exercism/cli/configuration"
	"os"
)

func Logout(dir string) {
	os.Remove(configuration.Filename(dir))
}
