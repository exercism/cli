package main

import (
	"os"
	"github.com/exercism/cli/configuration"
)

func Logout(dir string) {
	os.Remove(configuration.Filename(dir))
}
