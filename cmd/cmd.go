package cmd

import (
	"log"
	"os"
)

// Bail handles exitable errors in commands.
// TODO: figure out what goes here.
func Bail(err error) {
	log.Println(err)
	os.Exit(1)
}
