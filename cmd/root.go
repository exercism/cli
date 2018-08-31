package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	// BinaryName is the name of the app.
	// By default this is exercism, but people
	// are free to name this however they want.
	// The usage examples and help strings should reflect
	// the actual name of the binary.
	BinaryName string
	// Out is used to write to information.
	Out io.Writer
	// Err is used to write errors.
	Err io.Writer
	// In is used to provide mocked test input (i.e. for prompts).
	In io.Reader

	// TimeoutInSeconds is the timeout the default HTTP client will use.
	TimeoutInSeconds = 60

	// httpClient is the client used to make HTTP calls in the cli package.
	httpClient = &http.Client{Timeout: time.Duration(TimeoutInSeconds) * time.Second}
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   BinaryName,
	Short: "A friendly command-line interface to Exercism.",
	Long: `A command-line interface for the v2 redesign of Exercism.

Download exercises and submit your solutions.`,
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			api.Verbose = verbose
		}
		if timeout, _ := cmd.Flags().GetInt("timeout"); timeout > 0 {
			TimeoutInSeconds = timeout
		}

		httpClient = &http.Client{Timeout: time.Duration(TimeoutInSeconds) * time.Second}
		// share the configured HTTPClient with other web related services
		api.HTTPClient = httpClient
		cli.HTTPClient = httpClient
	},
}

// Execute adds all child commands to the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	BinaryName = os.Args[0]
	config.SetDefaultDirName(BinaryName)
	Out = os.Stdout
	Err = os.Stderr
	In = os.Stdin
	api.UserAgent = fmt.Sprintf("github.com/exercism/cli v%s (%s/%s)", Version, runtime.GOOS, runtime.GOARCH)
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().IntP("timeout", "", 0, "override the default HTTP timeout (seconds)")
}
