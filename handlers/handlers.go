package handlers

import (
	"path/filepath"
	"regexp"
)

const (
	msgPleaseAuthenticate = "You must be authenticated. Run `exercism configure --key=YOUR_API_KEY`."
)

func isTest(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".t" {
		return true
	}

	file := filepath.Base(path)
	name := file[:len(file)-len(ext)]
	if name == "test" {
		return true
	}
	return regexp.MustCompile(`[\._-]?[tT]est`).MatchString(name)
}
