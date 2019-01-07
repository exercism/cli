package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/browser"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// openCmd opens the designated exercise in the browser.
var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"o"},
	Short:   "Open an exercise on the website.",
	Long: `Open the specified exercise to the solution page on the Exercism website.

Pass the path to the directory that contains the solution you want to see on the website.
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
		cfg.UserViperConfig = v

		return runOpen(cfg, cmd.Flags(), args)
	},
}

func runOpen(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	var url string

	if remote, _ := flags.GetBool("remote"); remote {
		usrCfg := cfg.UserViperConfig

		apiUrl := fmt.Sprintf("%s/solutions/latest", usrCfg.GetString("apibaseurl"))

		client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
		if err != nil {
			return err
		}

		req, err := client.NewRequest("GET", apiUrl, nil)
		if err != nil {
			return err
		}

		q := req.URL.Query()
		q.Add("exercise_id", args[0])
		if track, _ := flags.GetString("track"); track != "" {
			q.Add("track_id", track)
		}
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		var payload openPayload
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			return fmt.Errorf("unable to parse API response - %s", err)
		}

		if res.StatusCode != http.StatusOK {
			switch payload.Error.Type {
			case "track_ambiguous":
				return fmt.Errorf("%s: %s", payload.Error.Message, strings.Join(payload.Error.PossibleTrackIDs, ", "))
			default:
				return errors.New(payload.Error.Message)
			}
		}

		url = payload.Solution.URL
	} else {
		metadata, err := workspace.NewExerciseMetadata(args[0])
		if err != nil {
			return err
		}
		url = metadata.URL
	}
	browser.Open(url)
	return nil
}

type openPayload struct {
	Solution struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"solution"`
	Error struct {
		Type             string   `json:"type"`
		Message          string   `json:"message"`
		PossibleTrackIDs []string `json:"possible_track_ids"`
	} `json:"error,omitempty"`
}

func setupOpenFlags(flags *pflag.FlagSet) {
	flags.BoolP("remote", "r", false, "checks for remote solutions")
	flags.StringP("track", "t", "", "the track id")
}

func init() {
	RootCmd.AddCommand(openCmd)
	setupOpenFlags(openCmd.Flags())
}
