package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"io"

	"github.com/exercism/cli/config"
	"github.com/spf13/viper"
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
)

const msgWelcomePleaseConfigure = `

    Welcome to Exercism!

    To get started, you need to configure the tool with your API token.
    Find your token at

        %s

    Then run the configure command:

        %s configure --token=YOUR_TOKEN

`

// Running configure without any arguments will attempt to
// set the default workspace. If the default workspace directory
// risks clobbering an existing directory, it will print an
// error message that explains how to proceed.
const msgRerunConfigure = `

    Please re-run the configure command to define where
    to download the exercises.

        %s configure
`

const msgMissingMetadata = `

    The exercise you are submitting doesn't have the necessary metadata.
    Please see https://github.com/exercism/website-copy/blob/main/pages/cli_v1_to_v2.md for instructions on how to fix it.

`

// validateUserConfig validates the presence of required user config values
func validateUserConfig(cfg *viper.Viper) error {
	if cfg.GetString("token") == "" {
		return fmt.Errorf(
			msgWelcomePleaseConfigure,
			config.SettingsURL(cfg.GetString("apibaseurl")),
			BinaryName,
		)
	}
	if cfg.GetString("workspace") == "" || cfg.GetString("apibaseurl") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}
	return nil
}

// decodedAPIError decodes and returns the error message from the API response.
// If the message is blank, it returns a fallback message with the status code.
func decodedAPIError(resp *http.Response) error {
	var apiError struct {
		Error struct {
			Type             string   `json:"type"`
			Message          string   `json:"message"`
			PossibleTrackIDs []string `json:"possible_track_ids"`
		} `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
		return fmt.Errorf("failed to parse API error response: %s", err)
	}
	if apiError.Error.Message != "" {
		if apiError.Error.Type == "track_ambiguous" {
			return fmt.Errorf(
				"%s: %s",
				apiError.Error.Message,
				strings.Join(apiError.Error.PossibleTrackIDs, ", "),
			)
		}
		return fmt.Errorf(apiError.Error.Message)
	}
	return fmt.Errorf("unexpected API response: %d", resp.StatusCode)
}
