package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
	app "github.com/urfave/cli"
)

// Upgrade allows the user to upgrade to the latest version of the CLI.
func Upgrade(ctx *app.Context) error {
	c := cli.New(ctx.App.Version)
	needed, err := c.IsUpgradeNeeded()
	if err != nil {
		return err
	}
	if needed {
		if err := c.Upgrade(); err != nil {
			return err
		}
	}
	fmt.Println("Your CLI is up to date!")
	return nil
}
