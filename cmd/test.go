package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Run the exercise's tests.",
	Long: `Run the exercise's tests.

	Run this command in an exercise's root directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTest(args)
	},
}

func runTest(args []string) error {
	track, err := getTrack()
	if err != nil {
		return err
	}

	testConf, ok := workspace.TestConfigurations[track]

	if !ok {
		return fmt.Errorf("the \"%s\" track does not yet support running tests using the Exercism CLI. Please see HELP.md for testing instructions", track)
	}

	command, err := testConf.GetTestCommand()
	if err != nil {
		return err
	}
	cmdParts := strings.Split(command, " ")

	// pass args/flags to this command down to the test handler
	if len(args) > 0 {
		cmdParts = append(cmdParts, args...)
	}

	fmt.Printf("Running tests via `%s`\n\n", strings.Join(cmdParts, " "))
	exerciseTestCmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	// pipe output directly out, preserving any color
	exerciseTestCmd.Stdout = os.Stdout
	exerciseTestCmd.Stderr = os.Stderr

	err = exerciseTestCmd.Run()
	if err != nil {
		// unclear what other errors would pop up here, but it pays to be defensive
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			// if subcommand returned a non-zero exit code, exit with the same
			os.Exit(exitCode)
		} else {
			log.Fatalf("Failed to get error from failed subcommand: %v", err)
		}
	}
	return nil
}

func getTrack() (string, error) {
	metadata, err := workspace.NewExerciseMetadata(".")
	if err != nil {
		return "", err
	}
	if metadata.Track == "" {
		return "", fmt.Errorf("no track found in exercise metadata")
	}

	return metadata.Track, nil
}

func init() {
	RootCmd.AddCommand(testCmd)
}
