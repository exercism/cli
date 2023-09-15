package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

type ExerciseSolution struct {
	Solution struct {
		ID   string `json:"id"`
		URL  string `json:"buildSolutionURL"`
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

type solutionFile struct {
	path      string
	sourceURL string
	slug      string
}

type solutionDownload struct {
	// either/or
	slug string
	uuid string

	// optional
	track          string
	team           string
	forceoverwrite bool

	solutionURL string

	solution *ExerciseSolution
}

func (sd *solutionDownload) set(flags *pflag.FlagSet) {
	var err error
	sd.uuid, err = flags.GetString("uuid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	sd.slug, err = flags.GetString("exercise")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	sd.track, err = flags.GetString("track")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	sd.team, err = flags.GetString("team")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}

	sd.forceoverwrite, err = flags.GetBool("force")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func (sd *solutionDownload) validate() error {
	var err error
	if err = sd.needsSlugXorUUID(); err != nil {
		return err
	}
	if err = sd.needsSlugWhenGivenTrackOrTeam(); err != nil {
		return err
	}
	return nil
}

// needsSlugXorUUID checks the presence of slug XOR uuid.
func (sd *solutionDownload) needsSlugXorUUID() error {
	if sd.slug != "" && sd.uuid != "" || sd.uuid == sd.slug {
		return errors.New("need an --exercise name or a solution --uuid")
	}
	return nil
}

// needsSlugWhenGivenTrackOrTeam ensures that track/team arguments are also given with a slug.
// (track/team meaningless when given a uuid).
func (sd *solutionDownload) needsSlugWhenGivenTrackOrTeam() error {
	if (sd.team != "" || sd.track != "") && sd.slug == "" {
		return errors.New("--track or --team requires --exercise (not --uuid)")
	}
	return nil
}

func (sd *solutionDownload) buildSolutionURL(apiURL string) {
	// buildSolutionURL
	id := "latest"
	if sd.uuid != "" {
		id = sd.uuid
	}
	sd.solutionURL = fmt.Sprintf("%s/solutions/%s", apiURL, id)

	// create new URL object
	url, err := netURL.Parse(sd.solutionURL)
	if err != nil {
		log.Fatal(err)
	}

	// buildQueryParams
	query := url.Query()
	if sd.slug != "" {
		query.Add("exercise_id", sd.slug)
		if sd.track != "" {
			query.Add("track_id", sd.track)
		}
		if sd.team != "" {
			query.Add("team_id", sd.team)
		}
	}
	url.RawQuery = query.Encode()

	sd.solutionURL = url.String()
}

func (sd *solutionDownload) fetchFiles(client *api.Client, usrCfg *viper.Viper) error {
	var err error

	metadata := sd.solution.GetSolutionMetadata()
	exerciseDir := metadata.Exercise(usrCfg.GetString("workspace")).MetadataDir()

	if _, err = os.Stat(exerciseDir); !sd.forceoverwrite && err == nil {
		return fmt.Errorf("directory '%s' already exists, use --force to overwrite", exerciseDir)
	}

	if err := os.MkdirAll(exerciseDir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := metadata.Write(exerciseDir); err != nil {
		return err
	}

	for _, sf := range sd.solution.files() {
		res, err := client.MakeRequest(sf.sourceURL, true)
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

func runDownload(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	token := usrCfg.GetString("token")
	apiURL := usrCfg.GetString("apibaseurl")

	client, err := api.NewClient(token, apiURL)
	if err != nil {
		return err
	}

	download, err := newDownload(client, flags, usrCfg)
	if err != nil {
		return err
	}

	if err := download.fetchFiles(client, usrCfg); err != nil {
		return err
	}
	return nil
}

func newDownload(client *api.Client, flags *pflag.FlagSet, usrCfg *viper.Viper) (*solutionDownload, error) {
	var err error

	apiURL := usrCfg.GetString("apibaseurl")

	d := &solutionDownload{}
	d.set(flags)
	if err = d.validate(); err != nil {
		return nil, err
	}
	d.buildSolutionURL(apiURL)

	res, err := client.MakeRequest(d.solutionURL, true)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, decodedAPIError(res)
	}

	body, _ := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewReader(body))

	if err := json.Unmarshal(body, &d.solution); err != nil {
		return nil, decodedAPIError(res)
	}

	return d, nil
}

func (es ExerciseSolution) GetSolutionMetadata() workspace.ExerciseMetadata {
	return workspace.ExerciseMetadata{
		AutoApprove:  es.Solution.Exercise.AutoApprove,
		Track:        es.Solution.Exercise.Track.ID,
		Team:         es.Solution.Team.Slug,
		ExerciseSlug: es.Solution.Exercise.ID,
		ID:           es.Solution.ID,
		URL:          es.Solution.URL,
		Handle:       es.Solution.User.Handle,
		IsRequester:  es.Solution.User.IsRequester,
	}
}

func (es ExerciseSolution) files() []solutionFile {
	files := make([]solutionFile, 0, len(es.Solution.Files))
	for _, file := range es.Solution.Files {
		sf := solutionFile{
			path:      file,
			sourceURL: fmt.Sprintf("%s%s", es.Solution.FileDownloadBaseURL, file),
			slug:      es.Solution.Exercise.ID,
		}
		files = append(files, sf)
	}
	return files
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

// downloadCmd represents the solutionDownload command
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
		cfg := LoadUserConfig()
		return runDownload(cfg, cmd.Flags(), args)
	},
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
