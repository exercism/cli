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

	cfg.Set("workspace", config.Resolve(viperConfig.GetString("workspace"), configuration.Home))

	if cfg.GetString("apibaseurl") == "" {
		cfg.Set("apibaseurl", configuration.DefaultBaseURL)
	}
	if cfg.GetString("workspace") == "" {
		cfg.Set("workspace", config.DefaultWorkspaceDir(configuration))
	}

	show, err := flags.GetBool("show")
	if err != nil {
		return err
	}
	if show {
		defer printCurrentConfig(configuration)
	}
	client, err := api.NewClient(cfg.GetString("token"), cfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	switch {
	case cfg.GetString("token") == "":
		fmt.Fprintln(Err, "There is no token configured, please set it using --token.")
	case flags.Lookup("token").Changed:
		// User set new token
		skipAuth, _ := flags.GetBool("skip-auth")
		if !skipAuth {
			ok, err := client.TokenIsValid()
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(Err, "The token is invalid.")
			}
		}
	default:
		// Validate existing token
		skipAuth, _ := flags.GetBool("skip-auth")
		if !skipAuth {
			ok, err := client.TokenIsValid()
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(Err, "The token is invalid.")
			}
			defer printCurrentConfig(configuration)
		}
	}
	return configuration.Save("user")
}

func printCurrentConfig(configuration config.Configuration) {
	w := tabwriter.NewWriter(Out, 0, 0, 2, ' ', 0)
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
	configureCmd.Flags().StringP("token", "t", "", "authentication token used to connect to the site")
	configureCmd.Flags().StringP("workspace", "w", "", "directory for exercism exercises")
	configureCmd.Flags().StringP("api", "a", "", "API base url")
	configureCmd.Flags().BoolP("show", "s", false, "show the current configuration")
	configureCmd.Flags().BoolP("skip-auth", "", false, "skip online token authorization check")

	viperConfig = viper.New()
	viperConfig.BindPFlag("token", configureCmd.Flags().Lookup("token"))
	viperConfig.BindPFlag("workspace", configureCmd.Flags().Lookup("workspace"))
	viperConfig.BindPFlag("apibaseurl", configureCmd.Flags().Lookup("api"))
}

func init() {
	RootCmd.AddCommand(configureCmd)

	initConfigureCmd()
}
