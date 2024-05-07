package cmd

import (
	"github.com/exercism/cli/browser"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
)

// openCmd opens the designated exercise in the browser.
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open an exercise on the website.",
	Long: `Open the specified exercise to the solution page on the Exercism website.

Pass the path to the directory that contains the solution you want to see on the website.
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}
		metadata, err := workspace.NewExerciseMetadata(path)
		if err != nil {
			return err
		}
		return browser.Open(metadata.URL)
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
