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
	ws "github.com/exercism/cli/workspace"
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

// download is a download from the Exercism API.
type download struct {
	params *downloadParams
	*downloadPayload
	downloadWriter
}

// newDownloadFromFlags initiates a download from flags.
// This is the primary interaction for downloading from the Exercism API.
func newDownloadFromFlags(usrCfg *viper.Viper, flags *pflag.FlagSet) (*download, error) {
	downloadParams, err := newDownloadParamsFromFlags(usrCfg, flags)
	if err != nil {
		return nil, err
	}
	return newDownload(downloadParams, &fileDownloadWriter{})
}

// newDownloadFromExercise initiates a download from an exercise.
// This is used to get metadata and isn't the primary interaction for downloading.
// Only allows writing metadata, not exercise files.
func newDownloadFromExercise(usrCfg *viper.Viper, exercise ws.Exercise) (*download, error) {
	downloadParams, err := newDownloadParamsFromExercise(usrCfg, exercise)
	if err != nil {
		return nil, err
	}
	return newDownload(downloadParams, &fileDownloadWriter{})
}

// newDownload initiates a download by requesting a downloadPayload from the Exercism API.
func newDownload(params *downloadParams, writer downloadWriter) (*download, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	d := &download{params: params}
	if err := d.setPayload(); err != nil {
		return nil, err
	}
	if err := d.setWriter(writer); err != nil {
		return nil, err
	}
	return d, d.validate()
}

// requestSolutionFile requests a Solution file from the API, returning an HTTP response.
// Non-200 responses and 0 Content-Length responses are swallowed, returning nil.
func (d download) requestSolutionFile(filename string) (*http.Response, error) {
	parsedURL, err := netURL.ParseRequestURI(
		fmt.Sprintf("%s%s", d.Solution.FileDownloadBaseURL, filename))
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(d.params.token, d.params.apibaseurl)
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

func (d *download) setWriter(writer downloadWriter) error {
	if writer == nil {
		errors.New("writer is empty")
	}
	if err := writer.init(d); err != nil {
		return err
	}
	d.downloadWriter = writer
	return nil
}

// setPayload sets the downloadPayload by getting a payload from the Exercism API.
func (d *download) setPayload() error {
	client, err := api.NewClient(d.params.token, d.params.apibaseurl)
	if err != nil {
		return err
	}

	req, err := client.NewRequest("GET", d.requestURL(), nil)
	if err != nil {
		return err
	}
	d.buildQuery(req.URL)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&d.downloadPayload); err != nil {
		return fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf(
			"unauthorized request. Please run the configure command. You can find your API token at %s/my/settings",
			config.InferSiteURL(d.params.apibaseurl),
		)
	}
	if res.StatusCode != http.StatusOK {
		switch d.Error.Type {
		case "track_ambiguous":
			return fmt.Errorf("%s: %s", d.Error.Message, strings.Join(d.Error.PossibleTrackIDs, ", "))
		default:
			return errors.New(d.Error.Message)
		}
	}
	return nil
}

func (d download) requestURL() string {
	id := "latest"
	if d.params.uuid != "" {
		id = d.params.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", d.params.apibaseurl, id)
}

func (d download) buildQuery(url *netURL.URL) {
	query := url.Query()
	if d.params.slug != "" {
		query.Add("exercise_id", d.params.slug)
		if d.params.track != "" {
			query.Add("track_id", d.params.track)
		}
		if d.params.team != "" {
			query.Add("team_id", d.params.team)
		}
	}
	url.RawQuery = query.Encode()
}

func (d download) metadata() ws.ExerciseMetadata {
	return ws.ExerciseMetadata{
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

func (d download) exercise() ws.Exercise {
	return ws.Exercise{
		Root:  d.solutionRoot(),
		Track: d.Solution.Exercise.Track.ID,
		Slug:  d.Solution.Exercise.ID,
	}
}

// solutionRoot builds the root path based on the solution
// being part of a team and/or owned by another user.
func (d download) solutionRoot() string {
	root := d.params.workspace

	if d.isTeamSolution() {
		root = filepath.Join(root, "teams", d.Solution.Team.Slug)
	}
	if d.solutionBelongsToOtherUser() {
		root = filepath.Join(root, "users", d.Solution.User.Handle)
	}
	return root
}

func (d download) isTeamSolution() bool {
	return d.Solution.Team.Slug != ""
}

func (d download) solutionBelongsToOtherUser() bool {
	return !d.Solution.User.IsRequester
}

// validate validates the presence of an ID and checks for an error message.
func (d download) validate() error {
	if d.Solution.ID == "" {
		return errors.New("download missing an ID")
	}
	if d.Error.Message != "" {
		return errors.New(d.Error.Message)
	}
	return nil
}

// downloadWriter writes download contents.
type downloadWriter interface {
	init(*download) error
	writeMetadata() error
	writeSolutionFiles() error
	destination() string
}

// fileDownloadWriter writes download contents to the file system.
type fileDownloadWriter struct {
	*download
}

func (d *fileDownloadWriter) init(dl *download) error {
	d.download = dl
	return dl.validate()
}

// writeMetadata writes the exercise metadata.
func (d fileDownloadWriter) writeMetadata() error {
	metadata := d.metadata()
	return metadata.Write(d.destination())
}

// writeSolutionFiles attempts to write each exercise file that is part of the downloaded Solution.
// An HTTP request is made using each filename and failed responses are swallowed.
// All successful file responses are written except when 0 Content-Length.
func (d fileDownloadWriter) writeSolutionFiles() error {
	if err := d.params.ensureWritable(); err != nil {
		return err
	}
	for _, filename := range d.Solution.Files {
		res, err := d.requestSolutionFile(filename)
		if err != nil {
			return err
		}
		if res == nil {
			continue
		}
		defer res.Body.Close()

		// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
		// TODO: handle --force flag to overwrite without asking.

		destination := filepath.Join(
			d.destination(),
			sanitizeLegacyFilepath(filename, d.exercise().Slug))
		if err = os.MkdirAll(filepath.Dir(destination), os.FileMode(0755)); err != nil {
			return err
		}
		f, err := os.Create(destination)
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

// destination is the download destination path.
func (d fileDownloadWriter) destination() string {
	return d.exercise().MetadataDir()
}

// downloadParams is required to create a download.
type downloadParams struct {
	// either/or
	slug, uuid string

	// user config
	token, apibaseurl, workspace string

	// optional
	track, team string

	// duck-type for downloadParams created from varying types
	downloadParamsFrom
}

func newDownloadParamsFromExercise(usrCfg *viper.Viper, exercise ws.Exercise) (*downloadParams, error) {
	d := &downloadParams{
		slug:               exercise.Slug,
		track:              exercise.Track,
		downloadParamsFrom: downloadParamsFromExercise{},
	}
	d.setFieldsFromConfig(usrCfg)
	return d, d.validate()
}

func newDownloadParamsFromFlags(usrCfg *viper.Viper, flags *pflag.FlagSet) (*downloadParams, error) {
	d := &downloadParams{downloadParamsFrom: downloadParamsFromFlags{}}
	d.setFieldsFromConfig(usrCfg)
	if err := d.setFieldsFromFlags(flags); err != nil {
		return nil, err
	}
	return d, d.validate()
}

// setFieldsFromFlags sets the fields derived from flags.
func (d *downloadParams) setFieldsFromFlags(flags *pflag.FlagSet) error {
	var err error
	d.uuid, err = flags.GetString("uuid")
	if err != nil {
		return err
	}
	d.slug, err = flags.GetString("exercise")
	if err != nil {
		return err
	}
	d.track, err = flags.GetString("track")
	if err != nil {
		return err
	}
	d.team, err = flags.GetString("team")
	if err != nil {
		return err
	}
	return nil
}

// setFieldsFromConfig sets the fields derived from the user config.
func (d *downloadParams) setFieldsFromConfig(usrCfg *viper.Viper) {
	d.token = usrCfg.GetString("token")
	d.apibaseurl = usrCfg.GetString("apibaseurl")
	d.workspace = usrCfg.GetString("workspace")
}

func (d *downloadParams) validate() error {
	validator := downloadParamsValidator{downloadParams: d}

	if err := validator.needsSlugXorUUID(); err != nil {
		return err
	}
	if err := validator.needsUserConfigValues(); err != nil {
		return err
	}
	if err := validator.needsSlugWhenGivenTrackOrTeam(); err != nil {
		return err
	}
	return nil
}

func (d downloadParams) ensureWritable() error {
	if !d.writeExerciseFilesPermitted() {
		return errors.New("writing exercise files not permitted")
	}
	return nil
}

// downloadParamsValidator contains validation rules for downloadParams.
type downloadParamsValidator struct {
	*downloadParams
}

// needsSlugXorUUID checks the presence of either a slug or a uuid (but not both).
func (d downloadParamsValidator) needsSlugXorUUID() error {
	if d.slug != "" && d.uuid != "" || d.uuid == d.slug {
		return d.missingSlugOrUUIDError()
	}
	return nil
}

// needsUserConfigValues checks the presence of required values from the user config.
func (d downloadParamsValidator) needsUserConfigValues() error {
	errMsg := "missing required user config: '%s'"
	if d.token == "" {
		return fmt.Errorf(errMsg, "token")
	}
	if d.apibaseurl == "" {
		return fmt.Errorf(errMsg, "apibaseurl")
	}
	if d.workspace == "" {
		return fmt.Errorf(errMsg, "workspace")
	}
	return nil
}

// needsSlugWhenGivenTrackOrTeam ensures that track/team arguments are also given with a slug.
// (track/team meaningless when given a uuid).
func (d downloadParamsValidator) needsSlugWhenGivenTrackOrTeam() error {
	if (d.team != "" || d.track != "") && d.slug == "" {
		return d.givenTrackOrTeamMissingSlugError()
	}
	return nil
}

// downloadParamsFrom is the interface to faciliate polymorphism when creating downloadParams from different types.
type downloadParamsFrom interface {
	writeExerciseFilesPermitted() bool
	missingSlugOrUUIDError() error
	givenTrackOrTeamMissingSlugError() error
}

type downloadParamsFromFlags struct{}

func (d downloadParamsFromFlags) writeExerciseFilesPermitted() bool { return true }

func (d downloadParamsFromFlags) missingSlugOrUUIDError() error {
	return errors.New("need an --exercise name or a solution --uuid")
}

func (d downloadParamsFromFlags) givenTrackOrTeamMissingSlugError() error {
	return errors.New("--track or --team requires --exercise (not --uuid)")
}

type downloadParamsFromExercise struct{}

func (d downloadParamsFromExercise) writeExerciseFilesPermitted() bool { return false }

func (d downloadParamsFromExercise) missingSlugOrUUIDError() error {
	return errors.New("need a 'slug' or a 'uuid'")
}

func (d downloadParamsFromExercise) givenTrackOrTeamMissingSlugError() error {
	return errors.New("programmer error - should never happen")
}

// sanitizeLegacyFilepath is a workaround for a path bug due to an early design
// decision (later reversed) to allow numeric suffixes for exercise directories,
// allowing people to have multiple parallel versions of an exercise.
func sanitizeLegacyFilepath(file, slug string) string {
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

// downloadPayload is an Exercism API response.
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
