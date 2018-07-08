package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/debug"
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
			debug.Verbose = verbose
		}
	},
}

// Execute adds all child commands to the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(Err, err.Error())
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
}
