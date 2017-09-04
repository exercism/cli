package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/exercism/cli/api"
	"github.com/spf13/cobra"
)

var (
	// BinaryName is the name of the app.
	// By default this is exercism, but people
	// are free to name this however they want.
	// The usage examples and help strings should reflect
	// the actual name of the binary.
	BinaryName string
	Out        io.Writer
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   BinaryName,
	Short: "A friendly command-line interface to Exercism.",
	Long: `A command-line interface for https://v2.exercism.io.

Download exercises and submit your solutions.`,
}

// Execute adds all child commands to the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	BinaryName = os.Args[0]
	Out = os.Stdout
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}
