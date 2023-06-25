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

type TestConfiguration struct {
	// The static portion of the test command, which will be run for every test on this track. Examples include `cargo test` or `go test`
	command string
	// Some tracks test by running a specific file, such as `ruby lasagna_test.rb`. Set this to `true` to look up and include the name of the default test file.
	AppendTestFiles bool
	// All args after `--` aren't parsed and are passed to the test command. Some languages (especially `rust`) expect an additional `--` between _their_ args. So instead of requiring a user to call `exercism test -- -- --include-ingored` to run all `rust` tests, set this to `true` to separate the args passed to the test runner by a `--` automatically.
	autoSeparateArgs bool
}

var testConfigurations = map[string]TestConfiguration{
	"go": {
		command: "go test",
	},
	"rust": {
		command:          "cargo test",
		autoSeparateArgs: true,
	},
	"ruby": {
		command:         "ruby",
		AppendTestFiles: true,
	},
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Infer and run the test command for an exercise.",
	Long: `Infer and run the test command for an exercise.

	Run this command in an exercise's root directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		track, err := getTrack()

		if err != nil {
			return err
		}

		testConf, ok := testConfigurations[track]

		if !ok {
			return fmt.Errorf("test handler for the `%s` track not yet implemented. Please see HELP.md for testing instructions", track)
		}

		cmdParts := strings.Split(testConf.command, " ")

		if testConf.AppendTestFiles {
			testFiles, err := getTestFiles()
			if err != nil {
				return err
			}
			cmdParts = append(cmdParts, testFiles...)
		}

		// pass args/flags to this command down to the test handler
		if len(args) > 0 {
			if testConf.autoSeparateArgs {
				cmdParts = append(cmdParts, "--")
			}
			cmdParts = append(cmdParts, args...)
		}

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
	},
}

func getTrack() (string, error) {
	metadata, err := workspace.NewExerciseMetadata(".")
	if err != nil {
		return "", err
	}

	return metadata.Track, nil
}

func getTestFiles() ([]string, error) {
	testFiles, err := workspace.NewExerciseConfig(".")
	if err != nil {
		return []string{}, err
	}
	return testFiles.Files.Test, nil
}

func init() {
	RootCmd.AddCommand(testCmd)
}
