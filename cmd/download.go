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

// ExerciseSolution is the container for the exercise solution API response.
type ExerciseSolution struct {
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

// solutionDownload is a helper container for managing the download process.
type solutionDownload struct {
	slug, uuid  string // slug and uuid are mutually exclusive and only one can be specified either one at run time
	track       string
	team        string
	solutionURL string

	solution *ExerciseSolution
}

// solutionFile is a helper container that holds the information needed to download the exercise files.
type solutionFile struct {
	path    string
	baseURL string
	slug    string
}

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
		cfg, err := LoadUserConfigFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid user config: %v\n", err)
			return err
		}
		return runDownload(cfg, cmd.Flags(), args)
	},
}

func runDownload(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if err := validateUserConfig(usrCfg); err != nil {
		return err
	}

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	download, err := newDownload(client, flags, usrCfg)
	if err != nil {
		return err
	}

	metadata := download.solution.metadata()
	exerciseDir := metadata.Exercise(usrCfg.GetString("workspace")).MetadataDir()

	forceDownload, err := flags.GetBool("force")
	if err != nil {
		return err
	}

	if err := createExerciseDir(exerciseDir, forceDownload); err != nil {
		return err
	}
	// This writes the metadata file to the exercise directory.
	if err := metadata.Write(exerciseDir); err != nil {
		return err
	}

	if err := download.getFiles(client, exerciseDir); err != nil {
		return err
	}

	fmt.Printf("\nDownloaded to %s\n", exerciseDir)
	return nil
}

// // createExerciseDir creates the exercise directory and checks if it already exists.
func createExerciseDir(dirName string, force bool) error {
	if _, err := os.Stat(dirName); !force && err == nil {
		return fmt.Errorf("directory '%s' already exists, use --force to overwrite", dirName)
	}
	if err := os.MkdirAll(dirName, os.FileMode(0755)); err != nil {
		return err
	}
	return nil
}

// newDownload creates a new solutionDownload object, which is a container for the exercise solution.
func newDownload(client *api.Client, flags *pflag.FlagSet, usrCfg *viper.Viper) (*solutionDownload, error) {
	var err error

	d := &solutionDownload{}
	d.set(flags)
	if err = d.validate(); err != nil {
		return nil, err
	}
	d.buildSolutionURL(usrCfg.GetString("apibaseurl"))

	res, err := client.MakeRequest("GET", d.solutionURL, true)
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

// getFiles wraps the main logic for preparing and downloading the exercise files.
func (sd *solutionDownload) getFiles(client *api.Client, exerciseDir string) error {
	for _, exerciseFile := range sd.solution.collectSolutionFiles() {
		if err := exerciseFile.fetchExerciseFiles(client, exerciseDir); err != nil {
			return err
		}
	}
	return nil
}

// set sets the flags for the solutionDownload object.
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
}

// validate validates the flags for the solutionDownload object
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

// buildSolutionURL builds the solution URL.
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
		panic(err)
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

// metadata returns the metadata for the solutionDownload object.
func (es ExerciseSolution) metadata() workspace.ExerciseMetadata {
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

// collectSolutionFiles returns a slice of solutionFile objects that are used to download the exercise files.
func (es ExerciseSolution) collectSolutionFiles() []solutionFile {
	files := make([]solutionFile, 0, len(es.Solution.Files))
	for _, file := range es.Solution.Files {
		sf := solutionFile{
			path:    file,
			baseURL: es.Solution.FileDownloadBaseURL,
			slug:    es.Solution.Exercise.ID,
		}
		files = append(files, sf)
	}
	return files
}

// fetchExerciseFiles downloads the exercise files and saves them to the exercise directory.
func (sf solutionFile) fetchExerciseFiles(client *api.Client, targetDir string) error {
	url, err := sf.createDownloadURL()
	if err != nil {
		return err
	}

	exerciseFilePath := sf.relativePath()

	res, err := client.MakeRequest("GET", url, true)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return // ignore
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"error downloading %s: %s",
			exerciseFilePath,
			res.Status,
		)
	}

	// Don't bother with empty files.
	if res.Header.Get("Content-Length") == "0" {
		return nil
	}

	targetExerciseDir := filepath.Join(targetDir, filepath.Dir(exerciseFilePath))
	if err = os.MkdirAll(targetExerciseDir, os.FileMode(0755)); err != nil {
		return err
	}

	exerciseFile, err := os.Create(filepath.Join(targetDir, exerciseFilePath))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return // ignore
		}
	}(exerciseFile)

	_, err = io.Copy(exerciseFile, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// createDownloadURL creates the download URL for the solutionFile object.
func (sf solutionFile) createDownloadURL() (string, error) {
	url, err := netURL.ParseRequestURI(fmt.Sprintf("%s%s", sf.baseURL, sf.path))

	if err != nil {
		return "", err
	}

	return url.String(), nil
}

// relativePath returns the relative path for the solutionFile object.
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
