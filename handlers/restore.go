package handlers

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

// Restore returns a user's solved problems.
func Restore(ctx *cli.Context) {
	c, err := config.Read(ctx.GlobalString("config"))
	if err != nil {
		fmt.Println(err)
		return
	}

	url := fmt.Sprintf("%s/api/v1/iterations/%s/restore", c.Hostname, c.APIKey)

	problems, err := api.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	hw := NewHomework(problems, c)
	err = hw.Save()
	if err != nil {
		fmt.Println(err)
		return
	}
	hw.Summarize()
}
