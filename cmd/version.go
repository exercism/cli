package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
	"github.com/spf13/cobra"
)

// Version is the version of the current build.
// It follows semantic versioning.
const Version = "3.5.4"

// checkLatest flag for version command.
var checkLatest bool

// versionCmd outputs the version of the CLI.
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Version outputs the version of CLI.",
	Long: `Version outputs the version of the exercism binary that is in use.

To check for the latest available version, call the command with the
--latest flag.
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(currentVersion())

		if checkLatest {
			c := cli.New(Version)
			l, err := checkForUpdate(c)
			if err != nil {
				return err
			}

			fmt.Println(l)
		}

		return nil
	},
}

// currentVersion returns a formatted version string for the Exercism CLI.
func currentVersion() string {
	return fmt.Sprintf("exercism version %s", Version)
}

// checkForUpdate verifies if the CLI is running the latest version.
// If the client is out of date, the function returns upgrade instructions.
func checkForUpdate(c *cli.CLI) (string, error) {

	ok, err := c.IsUpToDate()
	if err != nil {
		return "", err
	}

	if ok {
		return "Your CLI version is up to date.", nil
	}

	// Anything but ok is out of date.
	msg := fmt.Sprintf("A new CLI version is available. Run `exercism upgrade` to update to %s", c.LatestRelease.Version())
	return msg, nil

}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&checkLatest, "latest", "l", false, "check latest available version")
}
