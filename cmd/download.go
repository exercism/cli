package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"d"},
	Short:   "Download an exercise.",
	Long: `Download an exercise.

You may download an exercise to work on. If you've already
started working on it, the command will also download your
latest solution.

Download other people's solutions by providing the UUID.
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

		return runDownload(cfg, cmd.Flags(), args)
	},
}

func runDownload(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if usrCfg.GetString("token") == "" {
		return fmt.Errorf(msgWelcomePleaseConfigure, config.SettingsURL(usrCfg.GetString("apibaseurl")), BinaryName)
	}
	if usrCfg.GetString("workspace") == "" || usrCfg.GetString("apibaseurl") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}

	uuid, err := flags.GetString("uuid")
	if err != nil {
		return err
	}
	slug, err := flags.GetString("exercise")
	if err != nil {
		return err
	}
	if uuid != "" && slug != "" || uuid == slug {
		return errors.New("need an --exercise name or a solution --uuid")
	}

	track, err := flags.GetString("track")
	if err != nil {
		return err
	}

	team, err := flags.GetString("team")
	if err != nil {
		return err
	}

	urlParam := "latest"
	if uuid != "" {
		urlParam = uuid
	}

	params := downloadParams{
		cfg:      cfg,
		uuid:     uuid,
		slug:     slug,
		track:    track,
		team:     team,
		urlParam: urlParam,
	}
	payload, err := getDownloadPayload(params)
	if err != nil {
		return err
	}

	if err := writeMetadataFromPayload(payload, cfg); err != nil {
		return err
	}

	if err := writeSolutionFilesFromPayload(payload, cfg); err != nil {
		return err
	}

	fmt.Fprintf(Err, "\nDownloaded to\n%s\n", getExerciseDirFromPayload(payload, cfg))
	return nil
}

func getExerciseDirFromPayload(payload *downloadPayload, cfg config.Config) string {
	usrCfg := cfg.UserViperConfig

	root := usrCfg.GetString("workspace")
	if payload.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", payload.Solution.Team.Slug)
	}
	if !payload.Solution.User.IsRequester {
		root = filepath.Join(root, "users", payload.Solution.User.Handle)
	}
	exercise := workspace.Exercise{
		Root:  root,
		Track: payload.Solution.Exercise.Track.ID,
		Slug:  payload.Solution.Exercise.ID,
	}
	return exercise.MetadataDir()
}

func setupDownloadFlags(flags *pflag.FlagSet) {
	flags.StringP("uuid", "u", "", "the solution UUID")
	flags.StringP("track", "t", "", "the track ID")
	flags.StringP("exercise", "e", "", "the exercise slug")
	flags.StringP("team", "T", "", "the team slug")
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	setupDownloadFlags(downloadCmd.Flags())
}
