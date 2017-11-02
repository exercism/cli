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
	RunE: func(cmd *cobra.Command, args []string) error {
		usrCfg, err := config.NewUserConfig()
		if err != nil {
			return err
		}

		cliCfg, err := config.NewCLIConfig()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			args = []string{cwd}
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
				fmt.Println(err)
				continue
			}
			s, ok := option.(*workspace.Solution)
			if !ok {
				fmt.Println("something went wrong trying to pick that solution, not sure what happened")
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
				if info.IsDir() {
					return nil
				}
				ok, err := track.AcceptFilename(path)
				if err != nil {
					return err
				}
				if !ok {
					return nil
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

		for _, path := range paths {
			// Don't submit empty files
			info, err := os.Stat(path)
			if err != nil {
				return err
			}
			if info.Size() == 0 {
				fmt.Println("File %s was empty, skipping...", path)
				continue
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			filename := strings.Replace(path, filepath.Join(usrCfg.Workspace, solution.Track, solution.Exercise), "", -1)

			dirname := fmt.Sprintf("%s%s%s", string(os.PathSeparator), solution.Exercise, string(os.PathSeparator))
			pieces := strings.Split(path, dirname)
			filename = fmt.Sprintf("%s%s", string(os.PathSeparator), pieces[len(pieces)-1])

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

		apiCfg, err := config.NewAPIConfig()
		if err != nil {
			return err
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}
		req, err := client.NewRequest("PATCH", apiCfg.URL("submit", solution.ID), body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := client.Do(req, nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bb := &bytes.Buffer{}
		_, err = bb.ReadFrom(resp.Body)
		if err != nil {
			return err
		}

		fmt.Fprintf(Out, "Submitted. View at %s\n", solution.URL)
		return nil
	},
}

func initSubmitCmd() {
	// TODO
}

func init() {
	RootCmd.AddCommand(submitCmd)
}
