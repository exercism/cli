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
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
)

// submitCmd lets people upload a solution to the website.
var submitCmd = &cobra.Command{
	Use:     "submit",
	Aliases: []string{"s"},
	Short:   "Submit your solution to an exercise.",
	Long:    `Submit your solution to an Exercism exercise.`,
	Run: func(cmd *cobra.Command, args []string) {
		usrCfg, err := config.NewUserConfig()
		BailOnError(err)

		cliCfg, err := config.NewCLIConfig()
		BailOnError(err)

		if len(args) == 0 {
			args = []string{"."}
		}

		ws := workspace.New(usrCfg.Workspace)
		tx, err := workspace.NewTransmission(ws.Dir, args)
		BailOnError(err)

		dirs, err := ws.Locate(tx.Dir)
		BailOnError(err)

		sx, err := workspace.NewSolutions(dirs)
		BailOnError(err)

		var solution *workspace.Solution
		for {
			prompt := `
			We found more than one. Which one did you mean?
			Type the number of the one you want to select.

			%s
			> `
			s, err := sx.Pick(prompt)
			if err != nil {
				fmt.Println(err)
				continue
			}
			solution = s
			break
		}

		if !solution.IsRequester {
			BailOnError(errors.New("not your solution"))
		}
		track := cliCfg.Tracks[solution.Track]
		if track == nil {
			err := prepareTrack(solution.Track)
			BailOnError(err)
		}

		paths := tx.Files
		if len(paths) == 0 {
			walkFn := func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				ok, err := track.AcceptFilename(path)
				BailOnError(err)
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

		for _, path := range paths {
			file, err := os.Open(path)
			BailOnError(err)
			defer file.Close()

			filename := strings.Replace(path, filepath.Join(usrCfg.Workspace, solution.Track, solution.Exercise), "", -1)

			dirname := fmt.Sprintf("%s%s%s", string(os.PathSeparator), solution.Exercise, string(os.PathSeparator))
			pieces := strings.Split(path, dirname)
			filename = fmt.Sprintf("%s%s", string(os.PathSeparator), pieces[len(pieces)-1])

			part, err := writer.CreateFormFile("files[]", filename)
			BailOnError(err)
			_, err = io.Copy(part, file)
			BailOnError(err)
		}

		err = writer.Close()
		BailOnError(err)

		apiCfg, err := config.NewAPIConfig()
		BailOnError(err)

		client, err := api.NewClient()
		BailOnError(err)
		req, err := client.NewRequest("PATCH", fmt.Sprintf(apiCfg.URL("submit"), solution.ID), body)
		BailOnError(err)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := client.Do(req, nil)
		BailOnError(err)
		defer resp.Body.Close()

		bb := &bytes.Buffer{}
		_, err = bb.ReadFrom(resp.Body)
		BailOnError(err)
	},
}

func initSubmitCmd() {
	// todo
}

func init() {
	RootCmd.AddCommand(submitCmd)
}
