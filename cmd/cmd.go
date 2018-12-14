package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	netURL "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
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

// validateUserConfig validates the presense of required user config values
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

// downloadWriter writes metadata and Solution files from a downloadPayload to disk.
type downloadWriter struct {
	*downloadPayload
}

func newDownloadWriter(payload *downloadPayload) (*downloadWriter, error) {
	if err := payload.validate(); err != nil {
		return nil, err
	}
	return &downloadWriter{payload}, nil
}

func (d downloadWriter) writeMetadata() error {
	metadata := d.metadata()
	return metadata.Write(d.exercise().MetadataDir())
}

// writeSolutionFiles attempts to write each Solution file in the downloadPayload.
// An HTTP request is made for each file and failed responses are swallowed.
// All successful file responses are written except where empty.
func (d downloadWriter) writeSolutionFiles() error {
	for _, filename := range d.Solution.Files {
		res, err := d.requestFile(filename)
		if err != nil {
			return err
		}
		if res == nil {
			continue
		}
		defer res.Body.Close()

		// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
		// TODO: handle --force flag to overwrite without asking.

		sanitizedPath := d.sanitizeLegacyFilepath(filename, d.exercise().Slug)
		fileWritePath := filepath.Join(d.exercise().MetadataDir(), sanitizedPath)
		if err = os.MkdirAll(filepath.Dir(fileWritePath), os.FileMode(0755)); err != nil {
			return err
		}

		f, err := os.Create(fileWritePath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, res.Body); err != nil {
			return err
		}
	}
	return nil
}

// sanitizeLegacyFilepath is a workaround for a path bug due to an early design
// decision (later reversed) to allow numeric suffixes for exercise directories,
// allowing people to have multiple parallel versions of an exercise.
func (d downloadWriter) sanitizeLegacyFilepath(file, slug string) string {
	pattern := fmt.Sprintf(`\A.*[/\\]%s-\d*/`, slug)
	rgxNumericSuffix := regexp.MustCompile(pattern)
	if rgxNumericSuffix.MatchString(file) {
		file = string(rgxNumericSuffix.ReplaceAll([]byte(file), []byte("")))
	}
	// Rewrite paths submitted with an older, buggy client where the Windows
	// path is being treated as part of the filename.
	file = strings.Replace(file, "\\", "/", -1)
	return filepath.FromSlash(file)
}

// downloadParams is required to create a downloadPayload.
type downloadParams struct {
	usrCfg *viper.Viper
	uuid   string
	slug   string
	track  string
	team   string
}

func newDownloadParamsFromExercise(usrCfg *viper.Viper, exercise workspace.Exercise) (*downloadParams, error) {
	d := &downloadParams{usrCfg: usrCfg, slug: exercise.Slug, track: exercise.Track}
	return d, d.validate()
}

func newDownloadParamsFromFlags(usrCfg *viper.Viper, flags *pflag.FlagSet) (*downloadParams, error) {
	var err error
	d := &downloadParams{usrCfg: usrCfg}

	d.uuid, err = flags.GetString("uuid")
	if err != nil {
		return nil, err
	}
	d.slug, err = flags.GetString("exercise")
	if err != nil {
		return nil, err
	}

	if err = d.validate(); err != nil {
		return nil, errors.New("need an --exercise name or a solution --uuid")
	}

	d.track, err = flags.GetString("track")
	if err != nil {
		return nil, err
	}
	d.team, err = flags.GetString("team")
	if err != nil {
		return nil, err
	}
	return d, err
}

func (d *downloadParams) validate() error {
	if d.slug != "" && d.uuid != "" || d.uuid == d.slug {
		return errors.New("need a 'slug' or a 'uuid'")
	}
	return validateUserConfig(d.usrCfg)
}

type downloadPayload struct {
	*downloadParams
	Solution struct {
		ID   string `json:"id"`
		URL  string `json:"url"`
		Team struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"team"`
		User struct {
			Handle      string `json:"handle"`
			IsRequester bool   `json:"is_requester"`
		} `json:"user"`
		Exercise struct {
			ID              string `json:"id"`
			InstructionsURL string `json:"instructions_url"`
			AutoApprove     bool   `json:"auto_approve"`
			Track           struct {
				ID       string `json:"id"`
				Language string `json:"language"`
			} `json:"track"`
		} `json:"exercise"`
		FileDownloadBaseURL string   `json:"file_download_base_url"`
		Files               []string `json:"files"`
		Iteration           struct {
			SubmittedAt *string `json:"submitted_at"`
		}
	} `json:"solution"`
	Error struct {
		Type             string   `json:"type"`
		Message          string   `json:"message"`
		PossibleTrackIDs []string `json:"possible_track_ids"`
	} `json:"error,omitempty"`
}

// newDownloadPayload creates a payload by making an HTTP request to the API.
func newDownloadPayload(params *downloadParams) (*downloadPayload, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	d := &downloadPayload{downloadParams: params}

	client, err := api.NewClient(params.usrCfg.GetString("token"), params.usrCfg.GetString("apibaseurl"))
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", d.requestURL(), nil)
	if err != nil {
		return nil, err
	}
	d.buildQuery(req.URL)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		siteURL := config.InferSiteURL(params.usrCfg.GetString("apibaseurl"))
		return nil, fmt.Errorf("unauthorized request. Please run the configure command. You can find your API token at %s/my/settings", siteURL)
	}
	if res.StatusCode != http.StatusOK {
		switch d.Error.Type {
		case "track_ambiguous":
			return nil, fmt.Errorf("%s: %s", d.Error.Message, strings.Join(d.Error.PossibleTrackIDs, ", "))
		default:
			return nil, errors.New(d.Error.Message)
		}
	}
	return d, nil
}

func (d *downloadPayload) requestURL() string {
	id := "latest"
	if d.uuid != "" {
		id = d.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", d.usrCfg.GetString("apibaseurl"), id)
}

func (d *downloadPayload) buildQuery(url *netURL.URL) {
	query := url.Query()
	if d.uuid == "" {
		query.Add("exercise_id", d.slug)
		if d.track != "" {
			query.Add("track_id", d.track)
		}
		if d.team != "" {
			query.Add("team_id", d.team)
		}
	}
	url.RawQuery = query.Encode()
}

// requestFile requests a Solution file from the API, returning an HTTP response.
// Non 200 responses and zero length file responses are swallowed, returning nil.
func (d *downloadPayload) requestFile(filename string) (*http.Response, error) {
	if filename == "" {
		return nil, errors.New("filename is empty")
	}

	unparsedURL := fmt.Sprintf("%s%s", d.Solution.FileDownloadBaseURL, filename)
	parsedURL, err := netURL.ParseRequestURI(unparsedURL)
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(d.usrCfg.GetString("token"), d.usrCfg.GetString("apibaseurl"))
	req, err := client.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: deal with it
		return nil, nil
	}
	// Don't bother with empty files.
	if res.Header.Get("Content-Length") == "0" {
		return nil, nil
	}

	return res, nil
}

func (d *downloadPayload) exercise() workspace.Exercise {
	root := d.usrCfg.GetString("workspace")
	if d.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", d.Solution.Team.Slug)
	}
	if !d.Solution.User.IsRequester {
		root = filepath.Join(root, "users", d.Solution.User.Handle)
	}
	return workspace.Exercise{
		Root:  root,
		Track: d.Solution.Exercise.Track.ID,
		Slug:  d.Solution.Exercise.ID,
	}
}

func (d *downloadPayload) metadata() workspace.ExerciseMetadata {
	return workspace.ExerciseMetadata{
		AutoApprove: d.Solution.Exercise.AutoApprove,
		Track:       d.Solution.Exercise.Track.ID,
		Team:        d.Solution.Team.Slug,
		Exercise:    d.Solution.Exercise.ID,
		ID:          d.Solution.ID,
		URL:         d.Solution.URL,
		Handle:      d.Solution.User.Handle,
		IsRequester: d.Solution.User.IsRequester,
	}
}

func (d *downloadPayload) validate() error {
	if d.Solution.ID == "" {
		return errors.New("download payload is empty")
	}
	if d.Error.Message != "" {
		return errors.New(d.Error.Message)
	}
	return nil
}
