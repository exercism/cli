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

// downloadContext represents the required context around obtaining a Solution
// payload from the API and working with the its contents.
type downloadContext struct {
	*downloadParams
	payload *downloadPayload
}

// newDownloadContext creates a downloadContext, making an HTTP request
// to populate the payload.
func newDownloadContext(params *downloadParams) (*downloadContext, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	payload, err := newDownloadPayload(params)
	if err != nil {
		return nil, err
	}

	return &downloadContext{
		downloadParams: params,
		payload:        payload,
	}, nil
}

// writeSolutionFiles attempts to write each solution file in the payload.
// An HTTP request is made for each file and failed responses are swallowed.
// All successful file responses are written except where empty.
func (d *downloadContext) writeSolutionFiles(exercise workspace.Exercise) error {
	if err := d.payload.validate(); err != nil {
		return err
	}

	for _, filename := range d.payload.Solution.Files {
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

		relativePath := d.sanitizeLegacyFilepath(filename, exercise.Slug)
		dir := filepath.Join(exercise.MetadataDir(), filepath.Dir(relativePath))
		if err = os.MkdirAll(dir, os.FileMode(0755)); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(exercise.MetadataDir(), relativePath))
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

func (d *downloadContext) requestFile(filename string) (*http.Response, error) {
	if err := d.payload.validate(); err != nil {
		return nil, err
	}
	if filename == "" {
		return nil, errors.New("filename is empty")
	}

	unparsedURL := fmt.Sprintf("%s%s", d.payload.Solution.FileDownloadBaseURL, filename)
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

func (d *downloadContext) writeMetadata(exercise workspace.Exercise) error {
	metadata, err := d.metadata()
	if err != nil {
		return err
	}

	if err := metadata.Write(exercise.MetadataDir()); err != nil {
		return err
	}
	return nil
}

func (d *downloadContext) exercise() (workspace.Exercise, error) {
	if err := d.payload.validate(); err != nil {
		return workspace.Exercise{}, err
	}

	root := d.usrCfg.GetString("workspace")
	if d.payload.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", d.payload.Solution.Team.Slug)
	}
	if !d.payload.Solution.User.IsRequester {
		root = filepath.Join(root, "users", d.payload.Solution.User.Handle)
	}
	return workspace.Exercise{
		Root:  root,
		Track: d.payload.Solution.Exercise.Track.ID,
		Slug:  d.payload.Solution.Exercise.ID,
	}, nil
}

func (d *downloadContext) metadata() (workspace.ExerciseMetadata, error) {
	if err := d.payload.validate(); err != nil {
		return workspace.ExerciseMetadata{}, err
	}

	return workspace.ExerciseMetadata{
		AutoApprove: d.payload.Solution.Exercise.AutoApprove,
		Track:       d.payload.Solution.Exercise.Track.ID,
		Team:        d.payload.Solution.Team.Slug,
		Exercise:    d.payload.Solution.Exercise.ID,
		ID:          d.payload.Solution.ID,
		URL:         d.payload.Solution.URL,
		Handle:      d.payload.Solution.User.Handle,
		IsRequester: d.payload.Solution.User.IsRequester,
	}, nil
}

// sanitizeLegacyFilepath is a workaround for a path bug due to an early design
// decision (later reversed) to allow numeric suffixes for exercise directories,
// allowing people to have multiple parallel versions of an exercise.
func (d *downloadContext) sanitizeLegacyFilepath(file, slug string) string {
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

// downloadParams is required to create a downloadContext.
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
	return nil
}

type downloadPayload struct {
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

// newDownloadPayload gets a payload by making an HTTP request to the API.
func newDownloadPayload(params *downloadParams) (*downloadPayload, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	d := &downloadPayload{}

	client, err := api.NewClient(params.usrCfg.GetString("token"), params.usrCfg.GetString("apibaseurl"))
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", d.requestURL(params), nil)
	if err != nil {
		return nil, err
	}

	if err = d.buildQuery(params, req.URL); err != nil {
		return nil, err
	}
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

func (d *downloadPayload) requestURL(params *downloadParams) string {
	id := "latest"
	if params.uuid != "" {
		id = params.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", params.usrCfg.GetString("apibaseurl"), id)
}

func (d *downloadPayload) buildQuery(params *downloadParams, url *netURL.URL) error {
	if url == nil {
		return errors.New("url is empty")
	}

	query := url.Query()
	if params.uuid == "" {
		query.Add("exercise_id", params.slug)
		if params.track != "" {
			query.Add("track_id", params.track)
		}
		if params.team != "" {
			query.Add("team_id", params.team)
		}
	}
	url.RawQuery = query.Encode()

	return nil
}

func (d *downloadPayload) validate() error {
	if d == nil {
		return errors.New("download payload is empty")
	}
	if d.Error.Message != "" {
		return errors.New(d.Error.Message)
	}
	return nil
}
