package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// prepareCmd does necessary setup for Exercism and its tracks.
var prepareCmd = &cobra.Command{
	Use:     "prepare",
	Aliases: []string{"p"},
	Short:   "Prepare does generic setup for Exercism and its tracks.",
	Long: `Prepare downloads settings and dependencies for Exercism and the language tracks.

When called without any arguments, this downloads all the copy for the CLI so we
know what to say in all the various situations, as well as an up-to-date list
of the API endpoints to use.

When called with a track ID, it will download the files that the track maintainers
have said are necessary for the track in general. Any files that are only
necessary for a specific exercise will only be downloaded with the exercise.

To customize the CLI to suit your own preferences, use the configure command.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prepare called")
	},
}

func init() {
	RootCmd.AddCommand(prepareCmd)
}
