package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// workspaceCmd outputs the path to the person's workspace directory.
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Print out the path to your Exercism workspace.",
	Long: `Print out the path to your Exercism workspace.

This command can be used for scripting, or it can be combined with shell
commands to take you to your workspace.

For example you can run:

    cd $(exercism workspace)

On Windows, this will work only with Powershell, however you would
need to be on the same drive as your workspace directory. Otherwise
nothing will happen.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		return runWorkspace(v, cmd.Flags())
	},
}

func runWorkspace(v *viper.Viper, flags *pflag.FlagSet) error {
	if err := validateUserConfig(v); err != nil {
		return err
	}

	path := v.GetString("workspace")
	teamDir, _ := flags.GetString("team")
	if teamDir != "" {
		path = filepath.Join(path, "teams", teamDir)
	}

	ws, err := workspace.New(path)
	if err != nil {
		if os.IsNotExist(err) && teamDir != "" {
			return errors.New("team not found in workspace")
		}
		return err
	}

	if track, _ := flags.GetBool("track"); track {
		trackPaths, err := ws.TrackPaths()
		if err != nil {
			return err
		}

		for _, trackPath := range trackPaths {
			fmt.Fprintf(Out, "%s\n", filepath.Join(ws.Dir, trackPath))
		}

		return nil
	}

	if ex, _ := flags.GetBool("exercise"); ex {
		exercises, err := ws.Exercises()
		if err != nil {
			return err
		}

		for _, exercise := range exercises {
			fmt.Fprintf(Out, "%s\n", filepath.Join(ws.Dir, exercise.Path()))
		}

		return nil
	}

	fmt.Fprintf(Out, "%s\n", ws.Dir)
	return nil
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
	setupWorkspaceFlags(workspaceCmd.Flags())
}

func setupWorkspaceFlags(flags *pflag.FlagSet) {
	flags.BoolP("track", "t", false, "print track paths")
	flags.BoolP("exercise", "e", false, "print exercise paths")
	flags.StringP("team", "T", "", "the team slug")
}
