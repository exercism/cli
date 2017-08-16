package cmd

import (
	"fmt"

	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
)

// workspaceCmd outputs the path to the person's workspace directory.
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Output the path to your workspace.",
	Long: `Output the path to your workspace.

This command can be combined with shell commands to take you there.

For example you can run:

    cd $(exercism workspace)

On Windows, this command works in Powershell, however you
would need to be on the same drive as your workspace directory.
Otherwise nothing will happen.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		usrCfg, err := config.NewUserConfig()
		BailOnError(err)
		fmt.Println(usrCfg.Workspace)
	},
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
}
