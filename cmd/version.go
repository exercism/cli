package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the version of the current build.
// It follows semantic versioning.
const Version = "3.0.0-alpha.1"

// versionCmd outputs the version of the CLI.
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Version outputs the version of CLI.",
	Long: `Version outputs the version of the exercism binary that is in use.

To see if there is a more recent version available, call the command with the
--latest flag.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
