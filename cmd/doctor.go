package cmd

import (
	"fmt"

	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// doctorCmd reports workspace cleanup tasks.
var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"doc"},
	Short:   "Doctor reports workspace cleanup tasks.",
	Long: `Doctor reports workspace cleanup tasks.

	use --fixup to execute reported tasks.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		usrCfg := viper.New()
		usrCfg.AddConfigPath(cfg.Dir)
		usrCfg.SetConfigName("user")
		usrCfg.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = usrCfg.ReadInConfig()
		cfg.UserViperConfig = usrCfg

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("cli")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		return runDoctor(cfg, cmd.Flags(), args)
	},
}

func runDoctor(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if usrCfg.GetString("workspace") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}

	ws, err := workspace.New(usrCfg.GetString("workspace"))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	RootCmd.AddCommand(doctorCmd)
}
