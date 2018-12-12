package cmd

import (
	"errors"
	"fmt"

	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command.
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
	setup, err := newDownloadCmdSetup(cfg.UserViperConfig, flags)
	if err != nil {
		return err
	}

	ctx, err := newDownloadContext(cfg.UserViperConfig, setup.downloadArgs)
	if err != nil {
		return err
	}

	exercise, err := ctx.exercise()
	if err != nil {
		return err
	}

	if err = ctx.writeSolutionFiles(exercise); err != nil {
		return err
	}

	if err := ctx.writeMetadata(exercise); err != nil {
		return err
	}

	ctx.printResult(exercise)
	return nil
}

type downloadCmdSetup struct {
	downloadArgs map[string]string
}

func newDownloadCmdSetup(usrCfg *viper.Viper, flags *pflag.FlagSet) (*downloadCmdSetup, error) {
	setup := &downloadCmdSetup{}
	if err := setup.validateConfig(usrCfg); err != nil {
		return nil, err
	}
	if err := setup.populateDownloadArgs(flags); err != nil {
		return nil, err
	}
	return setup, nil
}

func (d *downloadCmdSetup) validateConfig(usrCfg *viper.Viper) error {
	if usrCfg.GetString("token") == "" {
		return fmt.Errorf(msgWelcomePleaseConfigure, config.SettingsURL(usrCfg.GetString("apibaseurl")), BinaryName)
	}
	if usrCfg.GetString("workspace") == "" || usrCfg.GetString("apibaseurl") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}
	return nil
}

func (d *downloadCmdSetup) populateDownloadArgs(flags *pflag.FlagSet) error {
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

	d.downloadArgs = map[string]string{
		"uuid":  uuid,
		"slug":  slug,
		"track": track,
		"team":  team,
	}

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
