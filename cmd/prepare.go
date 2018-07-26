package cmd

import "github.com/spf13/cobra"

// prepareCmd does necessary setup for Exercism and its tracks.
var prepareCmd = &cobra.Command{
	Use:     "prepare",
	Aliases: []string{"p"},
	Short:   "Prepare does setup for Exercism and its tracks.",
	Long: `Prepare downloads settings and dependencies for Exercism and the language tracks.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(prepareCmd)
}
