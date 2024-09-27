package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/debug"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   getCommandName(),
	Short: "A friendly command-line interface to Exercism.",
	Long: `A command-line interface for Exercism.

Download exercises and submit your solutions.`,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			debug.Verbose = verbose
		}
		if unmask, _ := cmd.Flags().GetBool("unmask-token"); unmask {
			debug.UnmaskAPIKey = unmask
		}
		if timeout, _ := cmd.Flags().GetInt("timeout"); timeout > 0 {
			cli.TimeoutInSeconds = timeout
			api.TimeoutInSeconds = timeout
		}
	},
}

// Execute adds all child commands to the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func getCommandName() string {
	return os.Args[0]
}

func init() {
	BinaryName = getCommandName()
	config.SetDefaultDirName(BinaryName)
	Out = os.Stdout
	Err = os.Stderr
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().IntP("timeout", "", 0, "override the default HTTP timeout (seconds)")
	RootCmd.PersistentFlags().BoolP("unmask-token", "", false, "will unmask the API during a request/response dump")
}
