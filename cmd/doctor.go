package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// doctorCmd reports workspace cleanup tasks.
var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"doc"},
	Short:   "Doctor reports workspace cleanup tasks.",
	Long: `Doctor reports workspace cleanup tasks.

	Use --fixup to run reported tasks.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		usrCfg := viper.New()
		usrCfg.AddConfigPath(cfg.Dir)
		usrCfg.SetConfigName("user")
		usrCfg.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = usrCfg.ReadInConfig()
		cfg.UserViperConfig = usrCfg

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("cli")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		return runDoctor(cfg, cmd.Flags())
	},
}

func runDoctor(cfg config.Config, flags *pflag.FlagSet) error {
	usrCfg := cfg.UserViperConfig
	if usrCfg.GetString("workspace") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}

	fixup, err := flags.GetBool("fixup")
	if err != nil {
		return err
	}
	if fixup {
		if err := runFixup(cfg); err != nil {
			return err
		}
	}

	return nil
}

func runFixup(cfg config.Config) error {
	usrCfg := cfg.UserViperConfig

	ws, err := workspace.New(usrCfg.GetString("workspace"))
	if err != nil {
		return err
	}

	exercises, err := ws.PotentialExercises()
	if err != nil {
		return err
	}

	if err := fixupMetadata(exercises, cfg); err != nil {
		return err
	}

	// TODO: add numeric suffix dir cleanup #699

	return nil
}

func fixupMetadata(exercises []workspace.Exercise, cfg config.Config) error {
	for _, exercise := range exercises {
		if _, err := exercise.MigrateLegacyMetadataFile(); err != nil {
			return err
		}
		if ok, _ := exercise.HasMetadata(); !ok {
			if err := downloadMetadata(exercise, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: extract this and cmd/download into download service?
// copied from download#runDownload
func downloadMetadata(exercise workspace.Exercise, cfg config.Config) error {
	usrCfg := cfg.UserViperConfig

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	param := "latest"
	url := fmt.Sprintf("%s/solutions/%s", usrCfg.GetString("apibaseurl"), param)

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("exercise_id", exercise.Slug)
	q.Add("track_id", exercise.Track)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var payload downloadPayload
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		siteURL := config.InferSiteURL(usrCfg.GetString("apibaseurl"))
		return fmt.Errorf("unauthorized request. Please run the configure command. You can find your API token at %s/my/settings", siteURL)
	}

	if res.StatusCode != http.StatusOK {
		switch payload.Error.Type {
		case "track_ambiguous":
			return fmt.Errorf("%s: %s", payload.Error.Message, strings.Join(payload.Error.PossibleTrackIDs, ", "))
		default:
			return errors.New(payload.Error.Message)
		}
	}

	metadata := workspace.ExerciseMetadata{
		AutoApprove: payload.Solution.Exercise.AutoApprove,
		Track:       payload.Solution.Exercise.Track.ID,
		Team:        payload.Solution.Team.Slug,
		Exercise:    payload.Solution.Exercise.ID,
		ID:          payload.Solution.ID,
		URL:         payload.Solution.URL,
		Handle:      payload.Solution.User.Handle,
		IsRequester: payload.Solution.User.IsRequester,
	}

	root := usrCfg.GetString("workspace")
	if metadata.Team != "" {
		root = filepath.Join(root, "teams", metadata.Team)
	}
	if !metadata.IsRequester {
		root = filepath.Join(root, "users", metadata.Handle)
	}

	err = metadata.Write(exercise.MetadataDir())
	if err != nil {
		return err
	}

	if err := res.Body.Close(); err != nil {
		return err
	}

	return nil
}

func setupDoctorFlags(flags *pflag.FlagSet) {
	flags.BoolP("fixup", "f", false, "run tasks")
	// TODO: --dry-run for fixup
}

func init() {
	RootCmd.AddCommand(doctorCmd)
	setupDoctorFlags(doctorCmd.Flags())
}
