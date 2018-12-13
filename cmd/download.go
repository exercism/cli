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
	if err := validateUserConfig(cfg.UserViperConfig); err != nil {
		return err
	}

	downloadParams, err := newDownloadParams(flags)
	if err != nil {
		return err
	}

	ctx, err := newDownloadContext(cfg.UserViperConfig, downloadParams)
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

	fmt.Fprintf(Err, "\nDownloaded to\n")
	fmt.Fprintf(Out, "%s\n", exercise.MetadataDir())

	return nil
}

// downloadParams represents the params for creating a downloadContext.
type downloadParams struct{}

// newDownloadParams returns a map of downloadParams populated from flags.
func newDownloadParams(flags *pflag.FlagSet) (map[string]string, error) {
	d := &downloadParams{}
	return d.get(flags)
}

func (d *downloadParams) get(flags *pflag.FlagSet) (map[string]string, error) {
	uuid, err := flags.GetString("uuid")
	if err != nil {
		return nil, err
	}
	slug, err := flags.GetString("exercise")
	if err != nil {
		return nil, err
	}

	if err = validateDownloadParams(d, uuid, slug); err != nil {
		return nil, err
	}

	track, err := flags.GetString("track")
	if err != nil {
		return nil, err
	}
	team, err := flags.GetString("team")
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"uuid":  uuid,
		"slug":  slug,
		"track": track,
		"team":  team,
	}, nil
}

func (d *downloadParams) downloadParamsError() error {
	return errors.New("need an --exercise name or a solution --uuid")
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
