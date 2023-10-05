package cmd

import (
	"bytes"
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

	metadata := download.payload.metadata()
	dir := metadata.Exercise(usrCfg.GetString("workspace")).MetadataDir()

	if _, err = os.Stat(dir); !download.forceoverwrite && err == nil {
		return fmt.Errorf("directory '%s' already exists, use --force to overwrite", dir)
	}

	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := metadata.Write(dir); err != nil {
		return err
	}

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	for _, sf := range download.payload.files() {
		url, err := sf.url()
		if err != nil {
			return err
		}

		req, err := client.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			// TODO: deal with it
			continue
		}
		// Don't bother with empty files.
		if res.Header.Get("Content-Length") == "0" {
			continue
		}

		path := sf.relativePath()
		dir := filepath.Join(metadata.Dir, filepath.Dir(path))
		if err = os.MkdirAll(dir, os.FileMode(0755)); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(metadata.Dir, path))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, res.Body)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(Err, "\nDownloaded to\n")
	fmt.Fprintf(Out, "%s\n", metadata.Dir)
	return nil
}

type download struct {
	// either/or
	slug, uuid string

	// user config
	token, apibaseurl, workspace string

	// optional
	track, team    string
	forceoverwrite bool

	payload *downloadPayload
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

	d.forceoverwrite, err = flags.GetBool("force")
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

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, decodedAPIError(res)
	}

	body, _ := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewReader(body))

	if err := json.Unmarshal(body, &d.payload); err != nil {
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

func (dp downloadPayload) metadata() workspace.ExerciseMetadata {
	return workspace.ExerciseMetadata{
		AutoApprove:  dp.Solution.Exercise.AutoApprove,
		Track:        dp.Solution.Exercise.Track.ID,
		Team:         dp.Solution.Team.Slug,
		ExerciseSlug: dp.Solution.Exercise.ID,
		ID:           dp.Solution.ID,
		URL:          dp.Solution.URL,
		Handle:       dp.Solution.User.Handle,
		IsRequester:  dp.Solution.User.IsRequester,
	}
}

func (dp downloadPayload) files() []solutionFile {
	fx := make([]solutionFile, 0, len(dp.Solution.Files))
	for _, file := range dp.Solution.Files {
		f := solutionFile{
			path:    file,
			baseURL: dp.Solution.FileDownloadBaseURL,
			slug:    dp.Solution.Exercise.ID,
		}
		fx = append(fx, f)
	}
	return fx
}

type solutionFile struct {
	path, baseURL, slug string
}

func (sf solutionFile) url() (string, error) {
	url, err := netURL.ParseRequestURI(fmt.Sprintf("%s%s", sf.baseURL, sf.path))

	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (sf solutionFile) relativePath() string {
	file := sf.path

	// Work around a path bug due to an early design decision (later reversed) to
	// allow numeric suffixes for exercise directories, letting people have
	// multiple parallel versions of an exercise.
	pattern := fmt.Sprintf(`\A.*[/\\]%s-\d*/`, sf.slug)
	rgxNumericSuffix := regexp.MustCompile(pattern)
	if rgxNumericSuffix.MatchString(sf.path) {
		file = string(rgxNumericSuffix.ReplaceAll([]byte(sf.path), []byte("")))
	}

	// Rewrite paths submitted with an older, buggy client where the Windows path is being treated as part of the filename.
	file = strings.Replace(file, "\\", "/", -1)

	return filepath.FromSlash(file)
}

func setupDownloadFlags(flags *pflag.FlagSet) {
	flags.StringP("uuid", "u", "", "the solution UUID")
	flags.StringP("track", "t", "", "the track ID")
	flags.StringP("exercise", "e", "", "the exercise slug")
	flags.StringP("team", "T", "", "the team slug")
	flags.BoolP("force", "F", false, "overwrite existing exercise directory")
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	setupDownloadFlags(downloadCmd.Flags())
}
