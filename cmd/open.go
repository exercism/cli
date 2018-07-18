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
		cfg := config.NewConfiguration()

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

		solutions, err := workspace.NewSolutions(paths)
		if err != nil {
			return err
		}

		if len(solutions) == 0 {
			return nil
		}

		if len(solutions) > 1 {
			var mine []*workspace.Solution
			for _, s := range solutions {
				if s.IsRequester {
					mine = append(mine, s)
				}
			}
			solutions = mine
		}

		selection := comms.NewSelection()
		for _, solution := range solutions {
			selection.Items = append(selection.Items, solution)
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
			solution, ok := option.(*workspace.Solution)
			if ok {
				browser.Open(solution.URL)
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
