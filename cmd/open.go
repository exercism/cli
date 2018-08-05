package cmd

import (
	"errors"
	"fmt"

	"github.com/exercism/cli/browser"
	"github.com/exercism/cli/comms"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// openCmd opens the designated exercise in the browser.
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open an exercise on the website.",
	Long: `Open the specified exercise to the solution page on the Exercism website.

Pass either the name of an exercise, or the path to the directory that contains
the solution you want to see on the website.
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		ws, err := workspace.New(v.GetString("workspace"))
		if err != nil {
			return err
		}

		paths, err := ws.Locate(args[0])
		if err != nil {
			return err
		}

		collection, err := workspace.NewMetadataCollection(paths)
		if err != nil {
			return err
		}

		if len(collection) == 0 {
			return nil
		}

		if len(collection) > 1 {
			var mine []*workspace.Metadata
			for _, metadata := range collection {
				if metadata.IsRequester {
					mine = append(mine, metadata)
				}
			}
			collection = mine
		}

		selection := comms.NewSelection()
		for _, metadata := range collection {
			selection.Items = append(selection.Items, metadata)
		}
		for {
			prompt := `
We found more than one. Which one did you mean?
Type the number of the one you want to select.

%s
> `
			option, err := selection.Pick(prompt)
			if err != nil {
				fmt.Println(err)
				continue
			}
			metadata, ok := option.(*workspace.Metadata)
			if ok {
				browser.Open(metadata.URL)
				return nil
			}
			if err != nil {
				return errors.New("should never happen")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
