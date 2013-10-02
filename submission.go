package main

import "strings"

var testExtensions = map[string]string{
	"ruby":    "_test.rb",
	"js":      ".spec.js",
	"elixir":  "_test.exs",
	"clojure": "_test.clj",
	"python":  "_test.py",
	"go":      "_test.go",
	"haskell": "_test.hs",
}

func IsTest(filename string) bool {
	for _, ext := range testExtensions {
		if strings.LastIndex(filename, ext) > 0 {
			return true
		}
	}
	return false
}
