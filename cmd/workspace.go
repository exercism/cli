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
	Short:   "Print out the path to your Exercism workspace.",
	Long: `Print out the path to your Exercism workspace.

This command can be used for scripting, or it can be combined with shell
commands to take you to your workspace.

For example you can run:

    cd $(exercism workspace)

On Windows, this will work only with Powershell, however you would
need to be on the same drive as your workspace directory. Otherwise
nothing will happen.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		usrCfg, err := config.NewUserConfig()
		if err != nil {
			return err
		}
		fmt.Println(usrCfg.Workspace)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
}
