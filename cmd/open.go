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

Open local or remote exercises by slug. Alternatively, you can pass a local exercise directory.
	`,
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
	usrCfg := cfg.UserViperConfig
	trackID, err := flags.GetString("track")
	if err != nil {
		return err
	}

	exerciseSlug, err := flags.GetString("exercise")
	if err != nil {
		return err
	}

	remote, err := flags.GetBool("remote")
	if err != nil {
		return err
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	if exerciseSlug == "" && path == "" {
		return fmt.Errorf("must provide an --exercise slug or an exercise path")
	}

	if exerciseSlug == "" {
		// if no --exercise is given, use original functionality
		metadata, err := workspace.NewExerciseMetadata(path)
		if err != nil {
			return err
		}
		return browser.Open(metadata.URL)
	}

	if remote {
		apiURL := fmt.Sprintf("%s/solutions/latest", usrCfg.GetString("apibaseurl"))

		client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
		if err != nil {
			return err
		}

		req, err := client.NewRequest("GET", apiURL, nil)
		if err != nil {
			return err
		}

		q := req.URL.Query()
		q.Add("exercise_id", exerciseSlug)
		if trackID != "" {
			q.Add("track_id", trackID)
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

		return browser.Open(url)
	}

	ws, err := workspace.New(usrCfg.GetString("workspace"))
	if err != nil {
		return err
	}
	exercises, err := ws.Exercises()
	if err != nil {
		return err
	}

	matchingExerciseMeta := make([]*workspace.ExerciseMetadata, 0, len(exercises))
	for _, exercise := range exercises {
		metaDir := exercise.MetadataDir()
		meta, err := workspace.NewExerciseMetadata(metaDir)
		if err != nil {
			return err
		}

		if meta.ExerciseSlug != exerciseSlug {
			continue
		}

		if trackID != "" && meta.Track != trackID {
			continue
		}

		matchingExerciseMeta = append(matchingExerciseMeta, meta)
	}

	switch len(matchingExerciseMeta) {
	case 0:
		return fmt.Errorf("No matching exercise found")
	case 1:
		url = matchingExerciseMeta[0].URL
		break
	default:
		tracks := make([]string, 0, len(matchingExerciseMeta))
		for _, exercise := range matchingExerciseMeta {
			tracks = append(tracks, exercise.Track)
		}

		return fmt.Errorf("Please specify a track ID: %s", strings.Join(tracks, ", "))
	}

	return browser.Open(url)
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
	flags.StringP("exercise", "e", "", "the exercise slug")
}

func init() {
	RootCmd.AddCommand(openCmd)
	setupOpenFlags(openCmd.Flags())
}
