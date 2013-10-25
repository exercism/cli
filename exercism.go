package main

import (
	"os"
	"github.com/msgehard/go-exercism/configuration"
)

func Logout(dir string) {
	os.Remove(configuration.Filename(dir))
}
