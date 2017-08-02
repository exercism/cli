package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/paths"
	app "github.com/urfave/cli"
)

// Debug provides information about the user's environment and configuration.
func Debug(ctx *app.Context) error {
	cli.HTTPClient = &http.Client{Timeout: 20 * time.Second}
	c := cli.New(ctx.App.Version)

	cfg, err := config.New(ctx.GlobalString("config"))
	if err != nil {
		return err
	}
	uc := config.UserConfig{
		Path:      cfg.File,
		Home:      paths.Home,
		Workspace: cfg.Dir,
		Token:     cfg.APIKey,
	}

	status := cli.NewStatus(c, uc)
	status.Censor = !ctx.Bool("full-api-key")
	s, err := status.Check()
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}
