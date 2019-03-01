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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"d"},
	Short:   "Download an exercise.",
	Long: `Download an exercise.

You may download an exercise to work on. If you've already
started working on it, the command will also download your
latest solution.

Download other people's solutions by providing the UUID.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()
		cfg.UserViperConfig = v

		return runDownload(cfg, cmd.Flags(), args)
	},
}

func runDownload(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if err := validateUserConfig(usrCfg); err != nil {
		return err
	}

	download, err := newDownload(flags, usrCfg)
	if err != nil {
		return err
	}

	writer := newDownloadWriter(download)
	if err := writer.writeSolutionFiles(); err != nil {
		return err
	}
	if err := writer.writeMetadata(); err != nil {
		return err
	}

	fmt.Fprintf(Err, "\nDownloaded to\n")
	fmt.Fprintf(Out, "%s\n", writer.destination())
	return nil
}

type download struct {
	// either/or
	slug, uuid string

	// user config
	token, apibaseurl, workspace string

	// optional
	track, team string

	*downloadPayload
}

func newDownload(flags *pflag.FlagSet, usrCfg *viper.Viper) (*download, error) {
	var err error
	d := &download{}
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

	d.token = usrCfg.GetString("token")
	d.apibaseurl = usrCfg.GetString("apibaseurl")
	d.workspace = usrCfg.GetString("workspace")

	if err = d.needsSlugXorUUID(); err != nil {
		return nil, err
	}
	if err = d.needsUserConfigValues(); err != nil {
		return nil, err
	}
	if err = d.needsSlugWhenGivenTrackOrTeam(); err != nil {
		return nil, err
	}

	client, err := api.NewClient(d.token, d.apibaseurl)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", d.url(), nil)
	if err != nil {
		return nil, err
	}
	d.buildQueryParams(req.URL)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&d.downloadPayload); err != nil {
		return nil, decodedAPIError(res)
	}

	return d, nil
}

func (d download) url() string {
	id := "latest"
	if d.uuid != "" {
		id = d.uuid
	}
	return fmt.Sprintf("%s/solutions/%s", d.apibaseurl, id)
}

func (d download) buildQueryParams(url *netURL.URL) {
	query := url.Query()
	if d.slug != "" {
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

// needsSlugXorUUID checks the presence of slug XOR uuid.
func (d download) needsSlugXorUUID() error {
	if d.slug != "" && d.uuid != "" || d.uuid == d.slug {
		return errors.New("need an --exercise name or a solution --uuid")
	}
	return nil
}

// needsUserConfigValues checks the presence of required values from the user config.
func (d download) needsUserConfigValues() error {
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
func (d download) needsSlugWhenGivenTrackOrTeam() error {
	if (d.team != "" || d.track != "") && d.slug == "" {
		return errors.New("--track or --team requires --exercise (not --uuid)")
	}
	return nil
}

func (d download) metadata() *ws.ExerciseMetadata {
	return &ws.ExerciseMetadata{
		AutoApprove:  d.Solution.Exercise.AutoApprove,
		Track:        d.Solution.Exercise.Track.ID,
		Team:         d.Solution.Team.Slug,
		ExerciseSlug: d.Solution.Exercise.ID,
		ID:           d.Solution.ID,
		URL:          d.Solution.URL,
		Handle:       d.Solution.User.Handle,
		IsRequester:  d.Solution.User.IsRequester,
	}
}

func (d download) exercise() ws.Exercise {
	return d.metadata().Exercise(d.workspace)
}

// requestFile requests a Solution file from the API, returning an HTTP response.
// 0 Content-Length responses are swallowed, returning nil.
func (d download) requestFile(filename string) (*http.Response, error) {
	parsedURL, err := netURL.ParseRequestURI(
		fmt.Sprintf("%s%s", d.Solution.FileDownloadBaseURL, filename))
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(d.token, d.apibaseurl)
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
	download *download
}

// newDownloadWriter creates a downloadWriter.
func newDownloadWriter(dl *download) *downloadWriter {
	return &downloadWriter{download: dl}
}

// writeMetadata writes the exercise metadata.
func (w downloadWriter) writeMetadata() error {
	return w.download.metadata().Write(w.destination())
}

// writeSolutionFiles attempts to write each file from the downloaded solution.
func (w downloadWriter) writeSolutionFiles() error {
	for _, filename := range w.download.Solution.Files {
		res, err := w.download.requestFile(filename)
		if err != nil {
			return err
		}
		if res == nil {
			// Ignore empty responses
			continue
		}
		defer res.Body.Close()

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
	}
	return nil
}

// destination is the download destination path.
func (w downloadWriter) destination() string {
	return w.download.exercise().MetadataDir()
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

func setupDownloadFlags(flags *pflag.FlagSet) {
	flags.StringP("uuid", "u", "", "the solution UUID")
	flags.StringP("track", "t", "", "the track ID")
	flags.StringP("exercise", "e", "", "the exercise slug")
	flags.StringP("team", "T", "", "the team slug")
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	setupDownloadFlags(downloadCmd.Flags())
}
