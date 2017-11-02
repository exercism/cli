package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viperUserConfig *viper.Viper
	viperAPIConfig  *viper.Viper
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
	Run: func(cmd *cobra.Command, args []string) {
		usrCfg := config.NewEmptyUserConfig()
		err := usrCfg.Load(viperUserConfig)
		BailOnError(err)

		apiCfg := config.NewEmptyAPIConfig()
		err = apiCfg.Load(viperAPIConfig)
		BailOnError(err)

		show, err := cmd.Flags().GetBool("show")
		BailOnError(err)

		if show {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()

			fmt.Fprintln(w, "")
			fmt.Fprintln(w, fmt.Sprintf("Config dir:\t%s", config.Dir()))
			fmt.Fprintln(w, fmt.Sprintf("-t, --token\t%s", usrCfg.Token))
			fmt.Fprintln(w, fmt.Sprintf("-w, --workspace\t%s", usrCfg.Workspace))
			fmt.Fprintln(w, fmt.Sprintf("-a, --api\t%s", apiCfg.BaseURL))
			return
		}

		err = usrCfg.Write()
		BailOnError(err)

		err = apiCfg.Write()
		BailOnError(err)
	},
}

func initConfigureCmd() {
	configureCmd.Flags().StringP("token", "t", "", "authentication token used to connect to exercism.io")
	configureCmd.Flags().StringP("workspace", "w", "", "directory for exercism exercises")
	configureCmd.Flags().StringP("api", "a", "", "API base url")
	configureCmd.Flags().BoolP("show", "s", false, "show the current configuration")

	viperUserConfig = viper.New()
	viperUserConfig.BindPFlag("token", configureCmd.Flags().Lookup("token"))
	viperUserConfig.BindPFlag("workspace", configureCmd.Flags().Lookup("workspace"))

	viperAPIConfig = viper.New()
	viperAPIConfig.BindPFlag("baseurl", configureCmd.Flags().Lookup("api"))
}

func init() {
	RootCmd.AddCommand(configureCmd)

	initConfigureCmd()
}
