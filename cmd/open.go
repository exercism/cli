package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// openCmd opens the designated exercise in the browser.
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open an exercise on the website.",
	Long: `Open the specified exercise to the solution page on the Exercism website.

Pass either the name of an exercise, or the path to the directory that contains
the solution you want to see on the website.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("open called")
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
