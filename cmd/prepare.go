package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// prepareCmd does necessary setup for Exercism and its tracks.
var prepareCmd = &cobra.Command{
	Use:     "prepare",
	Aliases: []string{"p"},
	Short:   "Prepare does setup for Exercism and its tracks.",
	Long: `Prepare downloads settings and dependencies for Exercism and the language tracks.

When called with a track ID, it will do specific setup for that track. This
might include downloading the files that the track maintainers have said are
necessary for the track in general. Any files that are only necessary for a specific
exercise will be downloaded along with the exercise.

To customize the CLI to suit your own preferences, use the configure command.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfiguration()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()
		cfg.UserViperConfig = v

		return runPrepare(cfg, cmd.Flags(), args)
	},
}

func runPrepare(cfg config.Configuration, flags *pflag.FlagSet, args []string) error {
	v := cfg.UserViperConfig

	track, err := flags.GetString("track")
	if err != nil {
		return err
	}

	if track == "" {
		return nil
	}
	client, err := api.NewClient(v.GetString("token"), v.GetString("apibaseurl"))
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/tracks/%s", v.GetString("apibaseurl"), track)

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var payload prepareTrackPayload

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(payload.Error.Message)
	}

	cliCfg, err := config.NewCLIConfig()
	if err != nil {
		return err
	}

	t, ok := cliCfg.Tracks[track]
	if !ok {
		t = config.NewTrack(track)
	}
	if payload.Track.TestPattern != "" {
		t.IgnorePatterns = append(t.IgnorePatterns, payload.Track.TestPattern)
	}
	cliCfg.Tracks[track] = t

	return cliCfg.Write()
}

type prepareTrackPayload struct {
	Track struct {
		ID          string `json:"id"`
		Language    string `json:"language"`
		TestPattern string `json:"test_pattern"`
	} `json:"track"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func initPrepareCmd() {
	prepareCmd.Flags().StringP("track", "t", "", "the track you want to prepare")
}

func init() {
	RootCmd.AddCommand(prepareCmd)
	initPrepareCmd()
}
