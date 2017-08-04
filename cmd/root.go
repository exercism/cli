package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/exercism/cli/api"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "exercism",
	Short: "A friendly command-line interface to Exercism.",
	Long: `A command-line interface for http://exercism.io.

Download exercises and submit your solution.`,
}

// Execute adds all child commands to the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
}
