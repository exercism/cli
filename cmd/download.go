package cmd

import (
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
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			return err
		}
		if token != "" {
			RootCmd.SetArgs([]string{"configure", "--token", token})
			if err := RootCmd.Execute(); err != nil {
				return err
			}
		}
		uuid, err := cmd.Flags().GetString("uuid")
		if err != nil {
			return err
		}
		if uuid == "" && len(args) == 0 {
			// TODO: usage
			return errors.New("need an exercise name or a solution --uuid")
		}
		usrCfg, err := config.NewUserConfig()
		if err != nil {
			return err
		}

		apiCfg, err := config.NewAPIConfig()
		if err != nil {
			return err
		}

		var slug string
		if uuid == "" {
			slug = "latest"
		} else {
			slug = uuid
		}
		url := fmt.Sprintf("%s/solutions/%s", apiCfg.BaseURL, slug)

		client, err := api.NewClient(usrCfg.Token, apiCfg.BaseURL)
		if err != nil {
			return err
		}

		req, err := client.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		track, err := cmd.Flags().GetString("track")
		if err != nil {
			return err
		}
		var exercise string
		if len(args) > 0 {
			exercise = args[0]
		}

		if uuid == "" {
			q := req.URL.Query()
			q.Add("exercise_id", exercise)
			if track != "" {
				q.Add("track_id", track)
			}
			req.URL.RawQuery = q.Encode()
		}

		payload := &downloadPayload{}
		res, err := client.Do(req, payload)
		if err != nil {
			return err
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

		var ws workspace.Workspace
		if solution.IsRequester {
			ws = workspace.New(filepath.Join(usrCfg.Workspace, solution.Track))
		} else {
			ws = workspace.New(filepath.Join(usrCfg.Workspace, "users", solution.Handle, solution.Track))
		}
		os.MkdirAll(ws.Dir, os.FileMode(0755))

		dir, err := ws.SolutionPath(solution.Exercise, solution.ID)
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

			res, err := client.Do(req, nil)
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
		fmt.Fprintf(Out, "\nDownloaded to\n%s\n", solution.Dir)
		return nil
	},
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

func initDownloadCmd() {
	downloadCmd.Flags().StringP("uuid", "u", "", "the solution UUID")
	downloadCmd.Flags().StringP("track", "t", "", "the track ID")
	downloadCmd.Flags().StringP("token", "k", "", "authentication token used to connect to the site")
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	initDownloadCmd()
}
