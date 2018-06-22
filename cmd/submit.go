package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/comms"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// submitCmd lets people upload a solution to the website.
var submitCmd = &cobra.Command{
	Use:     "submit",
	Aliases: []string{"s"},
	Short:   "Submit your solution to an exercise.",
	Long: `Submit your solution to an Exercism exercise.

The CLI will do its best to figure out what to submit.

If you call the command without any arguments, it will
submit the exercise contained in the current directory.

If called with the path to a directory, it will submit it.

If called with the name of an exercise, it will work out which
track it is on and submit it. The command will ask for help
figuring things out if necessary.
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check input before doing any other work
		exercise, err := cmd.Flags().GetString("exercise")
		if err != nil {
			return err
		}

		trackId, err := cmd.Flags().GetString("track")
		if err != nil {
			return err
		}

		files, err := cmd.Flags().GetStringSlice("files")
		if err != nil {
			return err
		}

		// Verify correct usage of flags/args
		if len(args) == 0 && len(files) == 0 && !(exercise != "" && trackId != "") {
			return errors.New("must use --exercise and --track together")
		}

		if len(args) > 0 && (exercise != "" || trackId != "" || len(files) > 0) {
			return errors.New("can't use flags and directory arguments together")
		}

		if len(files) > 0 && len(args) > 0 {
			return errors.New("can't submit files and a directory together")
		}

		usrCfg, err := config.NewUserConfig()
		if err != nil {
			return err
		}

		cliCfg, err := config.NewCLIConfig()
		if err != nil {
			return err
		}

		// TODO: make sure we get the workspace configured.
		if usrCfg.Workspace == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			usrCfg.Workspace = filepath.Dir(filepath.Dir(cwd))
		}

		ws := workspace.New(usrCfg.Workspace)

		// create directory from track and exercise slugs if needed
		if trackId != "" && exercise != "" {
			args = []string{filepath.Join(ws.Dir, trackId, exercise)}
		} else if len(files) > 0 {
			args = files
		}

		tx, err := workspace.NewTransmission(ws.Dir, args)
		if err != nil {
			return err
		}
		dirs, err := ws.Locate(tx.Dir)
		if err != nil {
			return err
		}
		sx, err := workspace.NewSolutions(dirs)
		if err != nil {
			return err
		}

		var solution *workspace.Solution

		selection := comms.NewSelection()
		for _, s := range sx {
			selection.Items = append(selection.Items, s)
		}

		for {
			prompt := `
			We found more than one. Which one did you mean?
			Type the number of the one you want to select.

			%s
			> `
			option, err := selection.Pick(prompt)
			if err != nil {
				return err
			}
			s, ok := option.(*workspace.Solution)
			if !ok {
				fmt.Fprintf(Err, "something went wrong trying to pick that solution, not sure what happened")
				continue
			}
			solution = s
			break
		}

		if !solution.IsRequester {
			return errors.New("not your solution")
		}

		track := cliCfg.Tracks[solution.Track]
		if track == nil {
			err := prepareTrack(solution.Track)
			if err != nil {
				return err
			}
			cliCfg.Load(viper.New())
			track = cliCfg.Tracks[solution.Track]
		}

		paths := tx.Files
		if len(paths) == 0 {
			walkFn := func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				ok, err := track.AcceptFilename(path)
				if err != nil || !ok {
					return err
				}
				paths = append(paths, path)
				return nil
			}
			filepath.Walk(solution.Dir, walkFn)
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		if len(paths) == 0 {
			return errors.New("no files found to submit")
		}

		// If the user submits a directory, confirm the list of files.
		if len(tx.ArgDirs) > 0 {
			prompt := "You specified a directory, which contains these files:\n"
			for i, path := range paths {
				prompt += fmt.Sprintf(" [%d]  %s\n", i+1, path)
			}
			prompt += "\nPress ENTER to submit, or control + c to cancel: "

			confirmQuestion := &comms.Question{
				Prompt:       prompt,
				DefaultValue: "y",
				Reader:       In,
				Writer:       Out,
			}
			answer, err := confirmQuestion.Ask()
			if err != nil {
				return err
			}
			if strings.ToLower(answer) != "y" {
				fmt.Fprintf(Err, "Submit cancelled.\nTry submitting individually instead.")
				return nil
			}
			fmt.Fprintf(Err, "Submitting files now...")
		}

		for _, path := range paths {
			// Don't submit empty files
			info, err := os.Stat(path)
			if err != nil {
				return err
			}
			if info.Size() == 0 {
				fmt.Printf("Warning: file %s was empty, skipping...", path)
				continue
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			dirname := fmt.Sprintf("%s%s%s", string(os.PathSeparator), solution.Exercise, string(os.PathSeparator))
			pieces := strings.Split(path, dirname)
			filename := fmt.Sprintf("%s%s", string(os.PathSeparator), pieces[len(pieces)-1])

			part, err := writer.CreateFormFile("files[]", filename)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
		}

		err = writer.Close()
		if err != nil {
			return err
		}

		client, err := api.NewClient(usrCfg.Token, usrCfg.APIBaseURL)
		if err != nil {
			return err
		}
		url := fmt.Sprintf("%s/solutions/%s", usrCfg.APIBaseURL, solution.ID)
		req, err := client.NewRequest("PATCH", url, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bb := &bytes.Buffer{}
		_, err = bb.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		if solution.AutoApprove == true {
			msg := `Your solution has been submitted successfully and has been auto-approved.
You can complete the exercise and unlock the next core exercise at:
`
			fmt.Fprintf(Err, msg)
		} else {
			msg := "Your solution has been submitted successfully. View it at:\n"
			fmt.Fprintf(Err, msg)
		}
		fmt.Fprintf(Out, "%s\n", solution.URL)

		return nil
	},
}

func initSubmitCmd() {
	submitCmd.Flags().StringP("track", "t", "", "the track ID")
	submitCmd.Flags().StringP("exercise", "e", "", "the exercise ID")
	submitCmd.Flags().StringSliceP("files", "f", make([]string, 0), "files to submit")
}

func init() {
	RootCmd.AddCommand(submitCmd)
}
