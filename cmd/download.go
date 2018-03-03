package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	Run: func(cmd *cobra.Command, args []string) {
		uuid, err := cmd.Flags().GetString("uuid")
		BailOnError(err)
		if uuid == "" && len(args) == 0 {
			// TODO: usage
			fmt.Fprintf(os.Stderr, "need an exercise name or a solution --uuid\n")
			return
		}
		apiCfg, err := config.NewAPIConfig()
		BailOnError(err)

		var slug string
		if uuid == "" {
			slug = "latest"
		} else {
			slug = uuid
		}
		url := fmt.Sprintf(apiCfg.URL("download"), slug)

		client, err := api.NewClient()
		BailOnError(err)

		req, err := client.NewRequest("GET", url, nil)
		BailOnError(err)

		track, err := cmd.Flags().GetString("track")
		BailOnError(err)
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
		BailOnError(err)

		if res.StatusCode != http.StatusOK {
			switch payload.Error.Type {
			case "track_ambiguous":
			default:
				fmt.Println(payload.Error.Message)
				os.Exit(1)
			}
		}

		solution := workspace.Solution{
			Track:       payload.Solution.Exercise.Track.ID,
			Exercise:    payload.Solution.Exercise.ID,
			ID:          payload.Solution.ID,
			URL:         payload.Solution.URL,
			Handle:      payload.Solution.User.Handle,
			IsRequester: payload.Solution.User.IsRequester,
		}

		var ws workspace.Workspace
		if solution.IsRequester {
			ws = workspace.New(filepath.Join(client.UserConfig.Workspace, solution.Track))
		} else {
			ws = workspace.New(filepath.Join(client.UserConfig.Workspace, "users", solution.Handle, solution.Track))
		}
		os.MkdirAll(ws.Dir, os.FileMode(0755))

		dir, err := ws.SolutionPath(solution.Exercise, solution.ID)
		BailOnError(err)

		os.MkdirAll(dir, os.FileMode(0755))

		err = solution.Write(dir)
		BailOnError(err)

		for _, file := range payload.Solution.Files {
			url := fmt.Sprintf("%s%s", payload.Solution.FileDownloadBaseURL, file)
			req, err := client.NewRequest("GET", url, nil)
			BailOnError(err)

			res, err := client.Do(req, nil)
			BailOnError(err)
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				// TODO: deal with it
				continue
			}

			// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
			// TODO: handle --force flag to overwrite without asking.
			relativePath := filepath.FromSlash(file)
			dir := filepath.Join(solution.Dir, filepath.Dir(relativePath))
			os.MkdirAll(dir, os.FileMode(0755))

			f, err := os.Create(filepath.Join(solution.Dir, relativePath))
			BailOnError(err)
			defer f.Close()
			_, err = io.Copy(f, res.Body)
			BailOnError(err)
		}
		fmt.Printf("\nDownloaded to\n%s\n", solution.Dir)
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
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	initDownloadCmd()
}
