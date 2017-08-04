package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd lets people download exercises and associated solutions.
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"d"},
	Short:   "Download an exercise.",
	Long: `Download an Exercism exercise to work on.

If you've already started working on the exercise, the command will
also download your most recent solution.

This does not automatically overwrite identically named files.
If it finds a file with the same name, it will ask if you want to overwrite it.
Or you can overwrite files automatically by passing the --force flag.

Download other people's solution by providing the UUID.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
