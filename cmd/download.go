package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		cfg := config.NewConfiguration()

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

func runDownload(cfg config.Configuration, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if usrCfg.GetString("token") == "" {
		return fmt.Errorf(msgWelcomePleaseConfigure, config.SettingsURL(usrCfg.GetString("apibaseurl")), BinaryName)
	}
	if usrCfg.GetString("workspace") == "" || usrCfg.GetString("apibaseurl") == "" {
		return fmt.Errorf(msgRerunConfigure, BinaryName)
	}

	uuid, err := flags.GetString("uuid")
	if err != nil {
		return err
	}
	exercise, err := flags.GetString("exercise")
	if err != nil {
		return err
	}
	if uuid == "" && exercise == "" {
		return errors.New("need an --exercise name or a solution --uuid")
	}

	var slug string
	if uuid == "" {
		slug = "latest"
	} else {
		slug = uuid
	}
	url := fmt.Sprintf("%s/solutions/%s", usrCfg.GetString("apibaseurl"), slug)

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	track, err := flags.GetString("track")
	if err != nil {
		return err
	}

	if uuid == "" {
		q := req.URL.Query()
		q.Add("exercise_id", exercise)
		if track != "" {
			q.Add("track_id", track)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var payload downloadPayload
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		siteURL := config.InferSiteURL(usrCfg.GetString("apibaseurl"))
		return fmt.Errorf("unauthorized request. Please run the configure command. You can find your API token at %s/my/settings", siteURL)
	}

	if res.StatusCode != http.StatusOK {
		switch payload.Error.Type {
		case "track_ambiguous":
			return fmt.Errorf("%s: %s", payload.Error.Message, strings.Join(payload.Error.PossibleTrackIDs, ", "))
		default:
			return errors.New(payload.Error.Message)
		}
	}

	solution := workspace.Solution{
		AutoApprove: payload.Solution.Exercise.AutoApprove,
		Track:       payload.Solution.Exercise.Track.ID,
		Exercise:    payload.Solution.Exercise.ID,
		ID:          payload.Solution.ID,
		URL:         payload.Solution.URL,
		Handle:      payload.Solution.User.Handle,
		IsRequester: payload.Solution.User.IsRequester,
	}

	dir := usrCfg.GetString("workspace")
	if !solution.IsRequester {
		dir = filepath.Join(dir, "users", solution.Handle)
	}
	dir = filepath.Join(dir, solution.Track)

	os.MkdirAll(dir, os.FileMode(0755))
	ws, err := workspace.New(dir)
	if err != nil {
		return err
	}

	dir, err = ws.SolutionPath(solution.Exercise, solution.ID)
	if err != nil {
		return err
	}

	os.MkdirAll(dir, os.FileMode(0755))

	err = solution.Write(dir)
	if err != nil {
		return err
	}

	for _, file := range payload.Solution.Files {
		url := fmt.Sprintf("%s%s", payload.Solution.FileDownloadBaseURL, file)
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

		// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
		// TODO: handle --force flag to overwrite without asking.
		relativePath := filepath.FromSlash(file)
		dir := filepath.Join(solution.Dir, filepath.Dir(relativePath))
		os.MkdirAll(dir, os.FileMode(0755))

		f, err := os.Create(filepath.Join(solution.Dir, relativePath))
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
	fmt.Fprintf(Out, "%s\n", solution.Dir)
	return nil
}

type downloadPayload struct {
	Solution struct {
		ID   string `json:"id"`
		URL  string `json:"url"`
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
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	setupDownloadFlags(downloadCmd.Flags())
}
