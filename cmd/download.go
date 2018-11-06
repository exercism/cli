package cmd

import (
	"errors"
	"fmt"

	"github.com/exercism/cli/config"
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
	if err := validateUserConfig(usrCfg); err != nil {
		return err
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

	var download = &downloadContext{
		usrCfg:  usrCfg,
		uuid:    uuid,
		slug:    slug,
		track:   track,
		team:    team,
	}
	if err := newDownload(download); err != nil {
		return err
	}
	if err := download.writeMetadata(); err != nil {
		return err
	}
	if err := download.writeSolutionFiles(); err != nil {
		return err
	}
	fmt.Fprintf(Err, "\nDownloaded to\n")
	fmt.Fprintf(Out, "%s\n", download.getExercise().MetadataDir())
	return nil
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
