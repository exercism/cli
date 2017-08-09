package cmd

import (
	"log"
	"os"
)

// BailOnError handles exitable errors in commands.
// TODO: figure out what goes here.
func BailOnError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
