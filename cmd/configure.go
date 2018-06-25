package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viperUserConfig *viper.Viper
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
		usrCfg := config.NewEmptyUserConfig()
		err := usrCfg.Load(viperUserConfig)
		if err != nil {
			return err
		}
		usrCfg.Workspace = config.Resolve(usrCfg.Workspace, usrCfg.Home)
		usrCfg.SetDefaults()

		show, err := cmd.Flags().GetBool("show")
		if err != nil {
			return err
		}
		if show {
			defer printCurrentConfig()
		}
		client, err := api.NewClient(usrCfg.Token, usrCfg.APIBaseURL)
		if err != nil {
			return err
		}

		switch {
		case usrCfg.Token == "":
			fmt.Fprintln(Out, "There is no token configured, please set it using --token.")
		case cmd.Flags().Lookup("token").Changed:
			// User set new token
			skipAuth, _ := cmd.Flags().GetBool("skip-auth")
			if !skipAuth {
				err = client.ValidateToken()
				if err != nil {
					return err
				}
			}
		default:
			// Validate existing token
			if !show {
				defer printCurrentConfig()
			}
			skipAuth, _ := cmd.Flags().GetBool("skip-auth")
			if !skipAuth {
				err = client.ValidateToken()
				if err != nil {
					fmt.Fprintln(Out, err)
				}
			}
		}

		return usrCfg.Write()
	},
}

func printCurrentConfig() {
	usrCfg, err := config.NewUserConfig()
	if err != nil {
		return
	}
	w := tabwriter.NewWriter(Out, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, fmt.Sprintf("Config dir:\t%s", config.Dir()))
	fmt.Fprintln(w, fmt.Sprintf("-t, --token\t%s", usrCfg.Token))
	fmt.Fprintln(w, fmt.Sprintf("-w, --workspace\t%s", usrCfg.Workspace))
	fmt.Fprintln(w, fmt.Sprintf("-a, --api\t%s", usrCfg.APIBaseURL))
	fmt.Fprintln(w, "")
}

func initConfigureCmd() {
	configureCmd.Flags().StringP("token", "t", "", "authentication token used to connect to the site")
	configureCmd.Flags().StringP("workspace", "w", "", "directory for exercism exercises")
	configureCmd.Flags().StringP("api", "a", "", "API base url")
	configureCmd.Flags().BoolP("show", "s", false, "show the current configuration")
	configureCmd.Flags().BoolP("skip-auth", "", false, "skip online token authorization check")

	viperUserConfig = viper.New()
	viperUserConfig.BindPFlag("token", configureCmd.Flags().Lookup("token"))
	viperUserConfig.BindPFlag("workspace", configureCmd.Flags().Lookup("workspace"))
	viperUserConfig.BindPFlag("apibaseurl", configureCmd.Flags().Lookup("api"))
}

func init() {
	RootCmd.AddCommand(configureCmd)

	initConfigureCmd()
}
