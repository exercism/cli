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

// download is a download from the Exercism API.
type download struct {
	params *downloadParams
	*downloadPayload
}

// newDownloadFromFlags initiates a download from flags.
// This is the primary interaction for downloading from the Exercism API.
func newDownloadFromFlags(flags *pflag.FlagSet, usrCfg *viper.Viper) (*download, error) {
	downloadParams, err := newDownloadParamsFromFlags(flags, usrCfg)
	if err != nil {
		return nil, err
	}
	return newDownload(downloadParams)
}

// newDownloadFromExercise initiates a download from an existing exercise directory.
// This is used to get metadata and isn't the primary interaction for downloading.
func newDownloadFromExercise(exercise ws.Exercise, usrCfg *viper.Viper) (*download, error) {
	downloadParams, err := newDownloadParamsFromExercise(exercise, usrCfg)
	if err != nil {
		return nil, err
	}
	return newDownload(downloadParams)
}

// newDownload creates a download by requesting a downloadPayload from the Exercism API.
func newDownload(params *downloadParams) (*download, error) {
	d := &download{params: params}

	client, err := api.NewClient(d.params.token, d.params.apibaseurl)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", d.payloadURL(), nil)
	if err != nil {
		return nil, err
	}
	d.buildPayloadQueryParams(req.URL)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&d.downloadPayload); err != nil {
		return nil, fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf(
			"unauthorized request. Please run the configure command. You can find your API token at %s/my/settings",
			config.InferSiteURL(d.params.apibaseurl),
		)
	}
	if res.StatusCode != http.StatusOK {
		switch d.Error.Type {
		case "track_ambiguous":
			return nil, fmt.Errorf(
				"%s: %s",
				d.Error.Message,
				strings.Join(d.Error.PossibleTrackIDs, ", "),
			)
		default:
			return nil, errors.New(d.Error.Message)
		}
	}

	return d, d.validate()
}

// payloadURL is the URL used to request a downloadPayload.
// The latest solution is used unless given a UUID.
func (d download) payloadURL() string {
	id := "latest"
	if d.params.uuid != "" {
		id = d.params.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", d.params.apibaseurl, id)
}

// buildPayloadQueryParams adds optional query parameters to the URL.
func (d download) buildPayloadQueryParams(url *netURL.URL) {
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
		Root:  d.solutionRootFilepath(),
		Track: d.Solution.Exercise.Track.ID,
		Slug:  d.Solution.Exercise.ID,
	}
}

// solutionRootFilepath builds the root path based on the solution
// being part of a team and/or owned by another user.
func (d download) solutionRootFilepath() string {
	root := d.params.workspace

	if d.isTeamSolution() {
		root = filepath.Join(root, "teams", d.Solution.Team.Slug)
	}
	if d.solutionBelongsToOtherUser() {
		root = filepath.Join(root, "users", d.Solution.User.Handle)
	}
	return root
}

// isTeamSolution indicates if the solution is part of a team.
func (d download) isTeamSolution() bool {
	return d.Solution.Team.Slug != ""
}

// solutionBelongsToOtherUser indicates if the solution belongs to another user
// (as opposed to being owned by the requesting user).
func (d download) solutionBelongsToOtherUser() bool {
	return !d.Solution.User.IsRequester
}

// ensureSolutionFilesWritable checks permission for writing solution files.
func (d download) ensureSolutionFilesWritable() error {
	if !d.params.downloadableFrom.writeSolutionFilesPermitted() {
		return errors.New("writing solution files not permitted in this context")
	}
	return nil
}

// validate verifies creation of a valid download.
func (d download) validate() error {
	if d.Solution.ID == "" {
		return errors.New("download missing an ID")
	}
	if d.Error.Message != "" {
		return errors.New(d.Error.Message)
	}
	return nil
}

// solutionFileRequester requests solution files from a download.
type solutionFileRequester struct {
	download *download
}

// request requests a Solution file from the API, returning an HTTP response.
// 0 Content-Length responses are swallowed, returning nil.
func (s solutionFileRequester) request(filename string) (*http.Response, error) {
	parsedURL, err := netURL.ParseRequestURI(
		fmt.Sprintf("%s%s", s.download.Solution.FileDownloadBaseURL, filename))
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(s.download.params.token, s.download.params.apibaseurl)
	req, err := client.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, decodedAPIError(res)
	}
	// Don't bother with empty files.
	if res.Header.Get("Content-Length") == "0" {
		return nil, nil
	}

	return res, nil
}

// downloadWriter writes download contents to the file system.
type downloadWriter struct {
	download      *download
	fileRequester solutionFileRequester
}

// newDownloadWriter creates a downloadWriter.
func newDownloadWriter(dl *download) *downloadWriter {
	return &downloadWriter{
		download:      dl,
		fileRequester: solutionFileRequester{download: dl},
	}
}

// writeMetadata writes the exercise metadata.
func (w downloadWriter) writeMetadata() error {
	metadata := w.download.metadata()
	return metadata.Write(w.destination())
}

func (w downloadWriter) writeSolutionFiles() error {
	for _, filename := range w.download.Solution.Files {
		if err := w.writeSolutionFile(filename); err != nil {
			return err
		}
	}
	return nil
}

// writeSolutionFile attempts to write the file from the downloaded solution.
func (w downloadWriter) writeSolutionFile(filename string) error {
	if err := w.download.ensureSolutionFilesWritable(); err != nil {
		return err
	}
	res, err := w.fileRequester.request(filename)
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}
	defer res.Body.Close()

	// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
	// TODO: handle --force flag to overwrite without asking.

	destination := filepath.Join(
		w.destination(),
		sanitizeLegacyFilepath(filename, w.download.exercise().Slug))
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
	return nil
}

// destination is the download destination path.
func (w downloadWriter) destination() string {
	return w.download.exercise().MetadataDir()
}

// downloadParams is parameter object for creating a download.
// A download may be constructed from multiple data sources; downloadParams encapsulates this construction.
type downloadParams struct {
	// either/or
	slug, uuid string

	// user config
	token, apibaseurl, workspace string

	// optional
	track, team string

	// polymorphic reference for downloads initiated via different types
	downloadableFrom
}

// newDownloadParamsFromExercise creates a new downloadParams from an exercise.
func newDownloadParamsFromExercise(exercise ws.Exercise, usrCfg *viper.Viper) (*downloadParams, error) {
	d := &downloadParams{
		slug:             exercise.Slug,
		track:            exercise.Track,
		downloadableFrom: downloadableFromExercise{},
	}
	return d.build(usrCfg)
}

// newDownloadParamsFromFlags creates a new downloadParams given flags.
func newDownloadParamsFromFlags(flags *pflag.FlagSet, usrCfg *viper.Viper) (*downloadParams, error) {
	d := &downloadParams{downloadableFrom: downloadableFromFlags{}}
	var err error
	d.uuid, err = flags.GetString("uuid")
	if err != nil {
		return nil, err
	}
	d.slug, err = flags.GetString("exercise")
	if err != nil {
		return nil, err
	}
	d.track, err = flags.GetString("track")
	if err != nil {
		return nil, err
	}
	d.team, err = flags.GetString("team")
	if err != nil {
		return nil, err
	}
	return d.build(usrCfg)
}

// build contains the common creation logic for creating downloadParams.
func (d *downloadParams) build(usrCfg *viper.Viper) (*downloadParams, error) {
	d.token = usrCfg.GetString("token")
	d.apibaseurl = usrCfg.GetString("apibaseurl")
	d.workspace = usrCfg.GetString("workspace")

	if d.downloadableFrom == nil {
		d.downloadableFrom = d
	}
	return d, d.validate()
}

// validate validates creation of downloadParams.
func (d downloadParams) validate() error {
	validator := downloadParamsValidator{params: &d}
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

// writeSolutionFilesPermitted is the default permission.
func (d downloadParams) writeSolutionFilesPermitted() bool { return false }

// errMissingSlugOrUUID is the default error.
func (d downloadParams) errMissingSlugOrUUID() error {
	return errors.New("programmer error")
}

// errGivenTrackOrTeamMissingSlug is the default error.
func (d downloadParams) errGivenTrackOrTeamMissingSlug() error {
	return errors.New("programmer error")
}

// downloadableFrom is a polymorphic interface for downloads initiated from different types.
// Clients can embed downloadParams to delegate to the default implementation.
type downloadableFrom interface {
	// indicates permission to write solution files from the payload.
	writeSolutionFilesPermitted() bool
	errMissingSlugOrUUID() error
	errGivenTrackOrTeamMissingSlug() error
}

// downloadableFromFlags represents downloadParams created from flags.
type downloadableFromFlags struct{}

func (d downloadableFromFlags) writeSolutionFilesPermitted() bool { return true }

func (d downloadableFromFlags) errMissingSlugOrUUID() error {
	return errors.New("need an --exercise name or a solution --uuid")
}

func (d downloadableFromFlags) errGivenTrackOrTeamMissingSlug() error {
	return errors.New("--track or --team requires --exercise (not --uuid)")
}

// downloadableFromExercise represents downloadParams created from an existing exercise directory.
// This delegates to the default implementation.
type downloadableFromExercise struct{ *downloadParams }

// downloadParamsValidator contains validation rules for downloadParams.
type downloadParamsValidator struct {
	params *downloadParams
}

// needsSlugXorUUID checks the presence of slug XOR uuid.
func (d downloadParamsValidator) needsSlugXorUUID() error {
	if d.params.slug != "" && d.params.uuid != "" || d.params.uuid == d.params.slug {
		return d.params.downloadableFrom.errMissingSlugOrUUID()
	}
	return nil
}

// needsUserConfigValues checks the presence of required values from the user config.
func (d downloadParamsValidator) needsUserConfigValues() error {
	errMsg := "missing required user config: '%s'"
	if d.params.token == "" {
		return fmt.Errorf(errMsg, "token")
	}
	if d.params.apibaseurl == "" {
		return fmt.Errorf(errMsg, "apibaseurl")
	}
	if d.params.workspace == "" {
		return fmt.Errorf(errMsg, "workspace")
	}
	return nil
}

// needsSlugWhenGivenTrackOrTeam ensures that track/team arguments are also given with a slug.
// (track/team meaningless when given a uuid).
func (d downloadParamsValidator) needsSlugWhenGivenTrackOrTeam() error {
	if (d.params.team != "" || d.params.track != "") && d.params.slug == "" {
		return d.params.downloadableFrom.errGivenTrackOrTeamMissingSlug()
	}
	return nil
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
