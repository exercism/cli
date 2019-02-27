package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"io"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/spf13/pflag"
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
    Please see https://exercism.io/cli-v1-to-v2 for instructions on how to fix it.

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
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
		return fmt.Errorf("failed to parse API error response: %s", err)
	}
	if apiError.Error.Message != "" {
		return fmt.Errorf(apiError.Error.Message)
	}
	return fmt.Errorf("unexpected API response: %d", resp.StatusCode)
}

type downloadClient struct {
	flags  *pflag.FlagSet
	usrCfg *viper.Viper
	uuid   string
	slug   string
	track  string
	team   string
}

func newDownloadClient(flags *pflag.FlagSet, usrCfg *viper.Viper) *downloadClient {
	return &downloadClient{
		flags:  flags,
		usrCfg: usrCfg,
	}
}

func (dlc *downloadClient) Do() (*http.Response, error) {
	if err := dlc.validateInput(); err != nil {
		return nil, err
	}

	client, err := api.NewClient(dlc.usrCfg.GetString("token"), dlc.usrCfg.GetString("apibaseurl"))
	if err != nil {
		return nil, err
	}

	url, err := dlc.url()
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = dlc.encodeQueryParams(req)

	return client.Do(req)
}

func (dlc *downloadClient) validateInput() error {
	uuid, err := dlc.flags.GetString("uuid")
	if err != nil {
		return err
	}
	dlc.uuid = uuid

	slug, err := dlc.flags.GetString("exercise")
	if err != nil {
		return err
	}
	dlc.slug = slug

	track, err := dlc.flags.GetString("track")
	if err != nil {
		return err
	}
	dlc.track = track

	team, err := dlc.flags.GetString("team")
	if err != nil {
		return err
	}
	dlc.team = team

	if uuid != "" && slug != "" || uuid == slug {
		return errors.New("need an --exercise name or a solution --uuid")
	}
	return nil
}

func (dlc downloadClient) url() (string, error) {
	identifier := "latest"
	if dlc.uuid != "" {
		identifier = dlc.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", dlc.usrCfg.GetString("apibaseurl"), identifier), nil
}

func (dlc *downloadClient) req(client *api.Client) (*http.Request, error) {
	url, err := dlc.url()
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = dlc.encodeQueryParams(req)
	return req, nil
}

func (dlc *downloadClient) encodeQueryParams(req *http.Request) *http.Request {
	if dlc.uuid != "" {
		return req
	}

	q := req.URL.Query()
	q.Add("exercise_id", dlc.slug)
	if dlc.track != "" {
		q.Add("track_id", dlc.track)
	}
	if dlc.team != "" {
		q.Add("team_id", dlc.team)
	}
	req.URL.RawQuery = q.Encode()
	return req
}
