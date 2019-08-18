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
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		metadata, err := workspace.NewExerciseMetadata(args[0])
		if err != nil {
			return err
		}
		return browser.Open(metadata.URL)
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
