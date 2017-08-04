package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// workspaceCmd outputs the path to the person's workspace directory.
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Output the path to your workspace.",
	Long: `Output the path to your workspace.

This command can be combined with shell commands to take you there.

For example, on Linux or MacOS you can run:

    cd $(exercism workspace)

On Windows, the equivalent command is:

    TODO ask @exercism/windows for help
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workspace called")
	},
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
}
