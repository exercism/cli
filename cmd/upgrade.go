package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
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
		quiet, _ := cmd.Flags().GetBool("quiet")

		c := cli.New(Version)
		return updateCLI(c, quiet)
	},
}

// updateCLI updates CLI to the latest available version, if it is out of date.
func updateCLI(c cli.Updater, quiet bool) error {
	ok, err := c.IsUpToDate()
	if err != nil {
		return err
	}

	if ok {
		if !quiet {
			fmt.Fprintln(Out, "Your CLI version is up to date.")
		}
		return nil
	}

	return c.Upgrade()
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
}
