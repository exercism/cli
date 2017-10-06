package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/debug"
	"github.com/spf13/cobra"
)

// upgradeCmd downloads and installs the most recent version of the CLI.
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"u"},
	Short:   "Upgrade to the latest version of the CLI.",
	Long: `Upgrade to the latest version of the CLI.

This finds and downloads the latest release, if you don't
already have it.

On Windows the old CLI will be left on disk, marked as hidden.
The next time you upgrade, the hidden file will be overwritten.
You can always delete this file.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			debug.Verbose = verbose
		}

		c := cli.New(Version)
		err := updateCLI(c)
		if err != nil {
			return err
		}
		return nil
	},
}

// updateCLI updates CLI to the latest available version, if it is out of date.
func updateCLI(c cli.Updater) error {
	ok, err := c.IsUpToDate()
	if err != nil {
		return err
	}

	if ok {
		fmt.Fprintln(Out, "Your CLI version is up to date.")
		return nil
	}

	return c.Upgrade()
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolP("verbose", "v", false, "verbose output")
}
