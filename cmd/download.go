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

	uuid, err := flags.GetString("uuid")
	if err != nil {
		return err
	}
	slug, err := flags.GetString("exercise")
	if err != nil {
		return err
	}
	if uuid != "" && slug != "" || uuid == slug {
		return errors.New("need an --exercise name or a solution --uuid")
	}

	track, err := flags.GetString("track")
	if err != nil {
		return err
	}

	team, err := flags.GetString("team")
	if err != nil {
		return err
	}

	urlParam := "latest"
	if uuid != "" {
		urlParam = uuid
	}

	params := downloadParams{
		cfg:      cfg,
		uuid:     uuid,
		slug:     slug,
		track:    track,
		team:     team,
		urlParam: urlParam,
	}
	payload, err := getDownloadPayload(params)
	if err != nil {
		return err
	}

	if err := writeMetadataFromPayload(payload, cfg); err != nil {
		return err
	}

	if err := writeSolutionFilesFromPayload(payload, cfg); err != nil {
		return err
	}

	fmt.Fprintf(Err, "\nDownloaded to\n%s\n", getExerciseDirFromPayload(payload, cfg))
	return nil
}

func getExerciseDirFromPayload(payload *downloadPayload, cfg config.Config) string {
	usrCfg := cfg.UserViperConfig

	root := usrCfg.GetString("workspace")
	if payload.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", payload.Solution.Team.Slug)
	}
	if !payload.Solution.User.IsRequester {
		root = filepath.Join(root, "users", payload.Solution.User.Handle)
	}
	exercise := workspace.Exercise{
		Root:  root,
		Track: payload.Solution.Exercise.Track.ID,
		Slug:  payload.Solution.Exercise.ID,
	}
	return exercise.MetadataDir()
}

type downloadParams struct {
	cfg      config.Config
	uuid     string
	slug     string
	track    string
	team     string
	urlParam string
}

func getDownloadPayload(params downloadParams) (*downloadPayload, error) {
	usrCfg := params.cfg.UserViperConfig

	url := fmt.Sprintf("%s/solutions/%s",
		usrCfg.GetString("apibaseurl"),
		params.urlParam,
	)

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if params.uuid == "" {
		q := req.URL.Query()
		q.Add("exercise_id", params.slug)
		if params.track != "" {
			q.Add("track_id", params.track)
		}
		if params.team != "" {
			q.Add("team_id", params.team)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var payload *downloadPayload
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		siteURL := config.InferSiteURL(usrCfg.GetString("apibaseurl"))
		return nil, fmt.Errorf("unauthorized request. Please run the configure command. You can find your API token at %s/my/settings", siteURL)
	}

	if res.StatusCode != http.StatusOK {
		switch payload.Error.Type {
		case "track_ambiguous":
			return nil, fmt.Errorf("%s: %s", payload.Error.Message, strings.Join(payload.Error.PossibleTrackIDs, ", "))
		default:
			return nil, errors.New(payload.Error.Message)
		}
	}

	return payload, nil
}

func writeMetadataFromPayload(payload *downloadPayload, cfg config.Config) error {
	if payload.Error.Message != "" {
		return errors.New(payload.Error.Message)
	}
	usrCfg := cfg.UserViperConfig

	metadata := workspace.ExerciseMetadata{
		AutoApprove: payload.Solution.Exercise.AutoApprove,
		Track:       payload.Solution.Exercise.Track.ID,
		Team:        payload.Solution.Team.Slug,
		Exercise:    payload.Solution.Exercise.ID,
		ID:          payload.Solution.ID,
		URL:         payload.Solution.URL,
		Handle:      payload.Solution.User.Handle,
		IsRequester: payload.Solution.User.IsRequester,
	}

	root := usrCfg.GetString("workspace")
	if metadata.Team != "" {
		root = filepath.Join(root, "teams", metadata.Team)
	}
	if !metadata.IsRequester {
		root = filepath.Join(root, "users", metadata.Handle)
	}

	exercise := workspace.Exercise{
		Root:  root,
		Track: metadata.Track,
		Slug:  metadata.Exercise,
	}

	dir := exercise.MetadataDir()

	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := metadata.Write(dir); err != nil {
		return err
	}

	return nil
}

func writeSolutionFilesFromPayload(payload *downloadPayload, cfg config.Config) error {
	if payload.Error.Message != "" {
		return errors.New(payload.Error.Message)
	}
	usrCfg := cfg.UserViperConfig

	root := usrCfg.GetString("workspace")
	if payload.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", payload.Solution.Team.Slug)
	}
	if !payload.Solution.User.IsRequester {
		root = filepath.Join(root, "users", payload.Solution.User.Handle)
	}
	exercise := workspace.Exercise{
		Root:  root,
		Track: payload.Solution.Exercise.Track.ID,
		Slug:  payload.Solution.Exercise.ID,
	}

	for _, file := range payload.Solution.Files {
		unparsedURL := fmt.Sprintf("%s%s", payload.Solution.FileDownloadBaseURL, file)
		parsedURL, err := netURL.ParseRequestURI(unparsedURL)
		if err != nil {
			return err
		}

		client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
		req, err := client.NewRequest("GET", parsedURL.String(), nil)
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

		// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
		// TODO: handle --force flag to overwrite without asking.

		// Work around a path bug due to an early design decision (later reversed) to
		// allow numeric suffixes for exercise directories, allowing people to have
		// multiple parallel versions of an exercise.
		// pattern := fmt.Sprintf(`\A.*[/\\]%s-\d*/`, metadata.Exercise)
		pattern := fmt.Sprintf(`\A.*[/\\]%s-\d*/`, payload.Solution.Exercise.ID)
		rgxNumericSuffix := regexp.MustCompile(pattern)
		if rgxNumericSuffix.MatchString(file) {
			file = string(rgxNumericSuffix.ReplaceAll([]byte(file), []byte("")))
		}

		// Rewrite paths submitted with an older, buggy client where the Windows path is being treated as part of the filename.
		file = strings.Replace(file, "\\", "/", -1)

		relativePath := filepath.FromSlash(file)

		dir := filepath.Join(exercise.MetadataDir(), filepath.Dir(relativePath))
		if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
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
