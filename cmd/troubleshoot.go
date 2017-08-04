package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// troubleshootCmd does a diagnostic self-check.
var troubleshootCmd = &cobra.Command{
	Use:     "troubleshoot",
	Aliases: []string{"t"},
	Short:   "Troubleshoot does a diagnostic self-check.",
	Long: `Troubleshoot provides output to help with troubleshooting.

If you're running into trouble, then run the troubleshoot command
and copy and paste the output into a GitHub issue so we can help
figure out what's going on.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("troubleshoot called")
	},
}

func init() {
	RootCmd.AddCommand(troubleshootCmd)
}
