package cmd

import (
	"fmt"
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
		configuration := config.NewConfiguration()

		viperConfig.AddConfigPath(configuration.Dir)
		viperConfig.SetConfigName("user")
		viperConfig.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = viperConfig.ReadInConfig()
		configuration.UserViperConfig = viperConfig

		return runConfigure(configuration, cmd.Flags())
	},
}

func runConfigure(configuration config.Configuration, flags *pflag.FlagSet) error {
	cfg := configuration.UserViperConfig

	show, err := flags.GetBool("show")
	if err != nil {
		return err
	}
	if show {
		printCurrentConfig(configuration)
		return nil
	}

	if flags.NFlag() == 0 && cfg.GetString("token") == "" {
		baseURL := cfg.GetString("apibaseurl")
		if baseURL != "" {
			tokenURL := config.InferSiteURL(baseURL) + "/my/settings"
			return fmt.Errorf("There is no token configured. Find your token on %s, and call this command again with --token=<your-token>.", tokenURL)
		}
		return fmt.Errorf("There is no token configured. Find your token in your settings on the website, and call this command again with --token=<your-token>.")
	}

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

	skipVerification, err := flags.GetBool("no-verify")
	if err != nil {
		return err
	}

	if !skipVerification {
		client, err := api.NewClient("", baseURL)
		if err != nil {
			return err
		}

		ok, err := client.IsPingable()
		if !ok || err != nil {
			return fmt.Errorf("The base API URL '%s' cannot be reached.\n\n%s", baseURL, err)
		}
	}
	cfg.Set("apibaseurl", baseURL)

	token, err := flags.GetString("token")
	if err != nil {
		return err
	}
	if token == "" {
		token = cfg.GetString("token")
	}

	tokenURL := config.InferSiteURL(cfg.GetString("apibaseurl")) + "/my/settings"
	if token == "" {
		return fmt.Errorf("There is no token configured. Find your token on %s, and call this command again with --token=<your-token>.", tokenURL)
	}

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
	cfg.Set("token", token)

	workspace, err := flags.GetString("workspace")
	if err != nil {
		return err
	}
	if workspace == "" {
		workspace = cfg.GetString("workspace")
	}
	workspace = config.Resolve(workspace, configuration.Home)
	if workspace == "" {
		workspace = config.DefaultWorkspaceDir(configuration)
	}
	cfg.Set("workspace", workspace)

	return configuration.Save("user")
}

func printCurrentConfig(configuration config.Configuration) {
	w := tabwriter.NewWriter(Err, 0, 0, 2, ' ', 0)
	defer w.Flush()

	v := configuration.UserViperConfig

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, fmt.Sprintf("Config dir:\t%s", configuration.Dir))
	fmt.Fprintln(w, fmt.Sprintf("-t, --token\t%s", v.GetString("token")))
	fmt.Fprintln(w, fmt.Sprintf("-w, --workspace\t%s", v.GetString("workspace")))
	fmt.Fprintln(w, fmt.Sprintf("-a, --api\t%s", v.GetString("apibaseurl")))
	fmt.Fprintln(w, "")
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
