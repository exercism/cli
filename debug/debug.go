package debug

import (
	"fmt"
	"io"
	"os"
)

var (
	// Verbose determines if debugging output is displayed to the user
	Verbose bool
	output  io.Writer = os.Stderr
)

// Println conditionally outputs a message to Stderr
func Println(args ...interface{}) {
	if Verbose {
		fmt.Fprintln(output, args...)
	}
}

// Printf conditionally outputs a formatted message to Stderr
func Printf(format string, args ...interface{}) {
	if Verbose {
		fmt.Fprintf(output, format, args...)
	}
}
