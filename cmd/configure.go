package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	viperConfig *viper.Viper
)

// configureCmd configures the command-line client with user-specific settings.
var configureCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"c"},
	Short:   "Configure the command-line client.",
	Long: `Configure the command-line client to customize it to your needs.

This lets you set up the CLI to talk to the API on your behalf,
and tells the CLI about your setup so it puts things in the right
places.

You can also override certain default settings to suit your preferences.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration := config.NewConfig()

		viperConfig.AddConfigPath(configuration.Dir)
		viperConfig.SetConfigName("user")
		viperConfig.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = viperConfig.ReadInConfig()
		configuration.UserViperConfig = viperConfig

		return runConfigure(configuration, cmd.Flags())
	},
}

func runConfigure(configuration config.Config, flags *pflag.FlagSet) error {
	cfg := configuration.UserViperConfig

	// Show the existing configuration and exit.
	show, err := flags.GetBool("show")
	if err != nil {
		return err
	}
	if show {
		printCurrentConfig(configuration)
		return nil
	}

	// If the command is run 'bare' and we have no token,
	// explain how to set the token.
	if flags.NFlag() == 0 && cfg.GetString("token") == "" {
		tokenURL := config.SettingsURL(cfg.GetString("apibaseurl"))
		return fmt.Errorf("There is no token configured. Find your token on %s, and call this command again with --token=<your-token>.", tokenURL)
	}

	// Determine the base API URL.
	baseURL, err := flags.GetString("api")
	if err != nil {
		return err
	}
	if baseURL == "" {
		baseURL = cfg.GetString("apibaseurl")
	}
	if baseURL == "" {
		baseURL = configuration.DefaultBaseURL
	}

	// By default we verify that
	// - the configured API URL is reachable.
	// - the configured token is valid.
	skipVerification, err := flags.GetBool("no-verify")
	if err != nil {
		return err
	}

	// Is the API URL reachable?
	if !skipVerification {
		client, err := api.NewClient("", baseURL)
		if err != nil {
			return err
		}

		if err := client.IsPingable(); err != nil {
			return fmt.Errorf("The base API URL '%s' cannot be reached.\n\n%s", baseURL, err)
		}
	}
	// Finally, configure the URL.
	cfg.Set("apibaseurl", baseURL)

	// Determine the token.
	token, err := flags.GetString("token")
	if err != nil {
		return err
	}
	if token == "" {
		token = cfg.GetString("token")
	}

	tokenURL := config.SettingsURL(cfg.GetString("apibaseurl"))

	// If we don't have a token then explain how to set it and bail.
	if token == "" {
		return fmt.Errorf("There is no token configured. Find your token on %s, and call this command again with --token=<your-token>.", tokenURL)
	}

	// Verify that the token is valid.
	if !skipVerification {
		client, err := api.NewClient(token, baseURL)
		if err != nil {
			return err
		}
		ok, err := client.TokenIsValid()
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("The token '%s' is invalid. Find your token on %s.", token, tokenURL)
		}
	}

	// Finally, configure the token.
	cfg.Set("token", token)

	// Determine the workspace.
	workspace, err := flags.GetString("workspace")
	if err != nil {
		return err
	}
	if workspace == "" {
		workspace = cfg.GetString("workspace")
	}
	workspace = config.Resolve(workspace, configuration.Home)

	if workspace != "" {
		// If there is a non-directory here, then we cannot proceed.
		if info, err := os.Lstat(workspace); !os.IsNotExist(err) && !info.IsDir() {
			msg := `

    There is already something at the workspace location you are configuring:

      %s

     Please rename it, or set a different workspace location:

       %s configure %s --workspace=PATH_TO_DIFFERENT_FOLDER
     `

			return fmt.Errorf(msg, workspace, BinaryName, commandify(flags))
		}
	}

	if workspace == "" {
		workspace = config.DefaultWorkspaceDir(configuration)

		// If it already exists don't clobber it with the default.
		if _, err := os.Lstat(workspace); !os.IsNotExist(err) {
			msg := `
    The default Exercism workspace is

      %s

    There is already something there.
    If it's a directory, that might be fine.
    If it's a file, you will need to move it first, or choose a
    different location for the workspace.

    You can choose the workspace location by rerunning this command
    with the --workspace flag.

      %s configure %s --workspace=%s
    `

			return fmt.Errorf(msg, workspace, BinaryName, commandify(flags), workspace)
		}
	}
	// Configure the workspace.
	cfg.Set("workspace", workspace)

	// Persist the new configuration.
	if err := configuration.Save("user"); err != nil {
		return err
	}
	fmt.Fprintln(Err, "\nYou have configured the Exercism command-line client:")
	printCurrentConfig(configuration)
	return nil
}

func printCurrentConfig(configuration config.Config) {
	w := tabwriter.NewWriter(Err, 0, 0, 2, ' ', 0)
	defer w.Flush()

	v := configuration.UserViperConfig

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, fmt.Sprintf("Config dir:\t\t%s", configuration.Dir))
	fmt.Fprintln(w, fmt.Sprintf("Token:\t(-t, --token)\t%s", v.GetString("token")))
	fmt.Fprintln(w, fmt.Sprintf("Workspace:\t(-w, --workspace)\t%s", v.GetString("workspace")))
	fmt.Fprintln(w, fmt.Sprintf("API Base URL:\t(-a, --api)\t%s", v.GetString("apibaseurl")))
	fmt.Fprintln(w, "")
}

func commandify(flags *pflag.FlagSet) string {
	var cmd string
	fn := func(f *pflag.Flag) {
		if f.Changed {
			cmd = fmt.Sprintf("%s --%s=%s", cmd, f.Name, f.Value.String())
		}
	}
	flags.VisitAll(fn)
	return strings.TrimLeft(cmd, " ")
}

func initConfigureCmd() {
	viperConfig = viper.New()
	setupConfigureFlags(configureCmd.Flags())
}

func setupConfigureFlags(flags *pflag.FlagSet) {
	flags.StringP("token", "t", "", "authentication token used to connect to the site")
	flags.StringP("workspace", "w", "", "directory for exercism exercises")
	flags.StringP("api", "a", "", "API base url")
	flags.BoolP("show", "s", false, "show the current configuration")
	flags.BoolP("no-verify", "", false, "skip online token authorization check")
}

func init() {
	RootCmd.AddCommand(configureCmd)

	initConfigureCmd()
}
