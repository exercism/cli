package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// submitCmd lets people upload a solution to the website.
var submitCmd = &cobra.Command{
	Use:     "submit FILE1 [FILE2 ...]",
	Aliases: []string{"s"},
	Short:   "Submit your solution to an exercise.",
	Long: `Submit your solution to an Exercism exercise.

    Call the command with the list of files you want to submit.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		usrCfg := viper.New()
		usrCfg.AddConfigPath(cfg.Dir)
		usrCfg.SetConfigName("user")
		usrCfg.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = usrCfg.ReadInConfig()
		cfg.UserViperConfig = usrCfg

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("cli")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		return runSubmit(cfg, cmd.Flags(), args)
	},
}

func runSubmit(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	if err := validateUserConfig(cfg.UserViperConfig); err != nil {
		return err
	}

	ctx, err := newSubmitCmdContext(cfg.UserViperConfig, flags, args)
	if err != nil {
		return err
	}

	if err := ctx.submit(cfg.UserViperConfig); err != nil {
		return err
	}

	ctx.printResult()
	return nil
}

// submitCmdContext represents the context for the submit cmd.
type submitCmdContext struct {
	usrCfg *viper.Viper
	flags  *pflag.FlagSet
	submission
}

// newSubmitCmdContext sets up a context to initiate a submission.
func newSubmitCmdContext(usrCfg *viper.Viper, flags *pflag.FlagSet, args []string) (*submitCmdContext, error) {
	ctx := &submitCmdContext{usrCfg: usrCfg, flags: flags}

	filepaths, err := ctx.sanitizeArgs(args)
	if err != nil {
		return nil, err
	}

	exercise, err := ctx.exercise(filepaths)
	if err != nil {
		return nil, err
	}

	if err = ctx.migrateLegacyMetadata(exercise); err != nil {
		return nil, err
	}

	documents, err := ctx._documents(filepaths, exercise)
	if err != nil {
		return nil, err
	}

	metadata, err := ctx._metadata(exercise)
	if err != nil {
		return nil, err
	}

	ctx.submission = submission{documents: documents, metadata: metadata}
	return ctx, nil
}

// sanitizeArgs validates args and swaps with evaluated symlink paths.
func (s *submitCmdContext) sanitizeArgs(args []string) ([]string, error) {
	for i, arg := range args {
		var err error
		arg, err = filepath.Abs(arg)
		if err != nil {
			return nil, err
		}

		info, err := os.Lstat(arg)
		if err != nil {
			if os.IsNotExist(err) {
				msg := `

    The file you are trying to submit cannot be found.

        %s

        `
				return nil, fmt.Errorf(msg, arg)
			}
			return nil, err
		}
		if info.IsDir() {
			msg := `

    You are submitting a directory, which is not currently supported.

        %s

    Please change into the directory and provide the path to the file(s) you wish to submit

        %s submit FILENAME

            `
			return nil, fmt.Errorf(msg, arg, BinaryName)
		}

		src, err := filepath.EvalSymlinks(arg)
		if err != nil {
			return nil, err
		}
		args[i] = src
	}
	return args, nil
}

func (s *submitCmdContext) exercise(filepaths []string) (workspace.Exercise, error) {
	ws, err := workspace.New(s.usrCfg.GetString("workspace"))
	if err != nil {
		return workspace.Exercise{}, err
	}

	var exerciseDir string
	for _, f := range filepaths {
		dir, err := ws.ExerciseDir(f)
		if err != nil {
			if workspace.IsMissingMetadata(err) {
				return workspace.Exercise{}, errors.New(msgMissingMetadata)
			}
			return workspace.Exercise{}, err
		}
		if exerciseDir != "" && dir != exerciseDir {
			msg := `

    You are submitting files belonging to different solutions.
    Please submit the files for one solution at a time.

        `
			return workspace.Exercise{}, errors.New(msg)
		}
		exerciseDir = dir
	}

	return workspace.NewExerciseFromDir(exerciseDir), nil
}

func (s *submitCmdContext) migrateLegacyMetadata(exercise workspace.Exercise) error {
	migrationStatus, err := exercise.MigrateLegacyMetadataFile()
	if err != nil {
		return err
	}
	if verbose, _ := s.flags.GetBool("verbose"); verbose {
		fmt.Fprintf(Err, migrationStatus.String())
	}
	return nil
}

func (s *submitCmdContext) _metadata(exercise workspace.Exercise) (*workspace.ExerciseMetadata, error) {
	metadata, err := workspace.NewExerciseMetadata(exercise.Filepath())
	if err != nil {
		return nil, err
	}

	if metadata.Exercise != exercise.Slug {
		// TODO: error msg should suggest running future doctor command
		msg := `

	The exercise directory does not match exercise slug in metadata:

		expected '%[1]s' but got '%[2]s'

	Please rename the directory '%[1]s' to '%[2]s' and try again.

		`
		return nil, fmt.Errorf(msg, exercise.Slug, metadata.Exercise)
	}

	if !metadata.IsRequester {
		// TODO: add test
		msg := `

    The solution you are submitting is not connected to your account.
    Please re-download the exercise to make sure it has the data it needs.

        %s download --exercise=%s --track=%s

        `
		return nil, fmt.Errorf(msg, BinaryName, metadata.Exercise, metadata.Track)
	}
	return metadata, nil
}

func (s *submitCmdContext) _documents(filepaths []string, exercise workspace.Exercise) ([]workspace.Document, error) {
	docs := make([]workspace.Document, 0, len(filepaths))
	for _, file := range filepaths {
		// Don't submit empty files
		info, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		const maxFileSize int64 = 65535
		if info.Size() >= maxFileSize {
			msg := `

      The submitted file '%s' is larger than the max allowed file size of %d bytes.
      Please reduce the size of the file and try again.

            `
			return nil, fmt.Errorf(msg, file, maxFileSize)
		}
		if info.Size() == 0 {

			msg := `

    WARNING: Skipping empty file
             %s

        `
			fmt.Fprintf(Err, msg, file)
			continue
		}
		doc, err := workspace.NewDocument(exercise.Filepath(), file)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if len(docs) == 0 {
		msg := `

    No files found to submit.

        `
		return nil, errors.New(msg)
	}
	return docs, nil
}

type submission struct {
	documents []workspace.Document
	metadata  *workspace.ExerciseMetadata
}

// submit submits the documents to the Exercism API.
func (s submission) submit(usrCfg *viper.Viper) error {
	if err := s.validate(); err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, doc := range s.documents {
		file, err := os.Open(doc.Filepath())
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("files[]", doc.Path())
		if err != nil {
			return err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/solutions/%s", usrCfg.GetString("apibaseurl"), s.metadata.ID)
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

	if resp.StatusCode == http.StatusBadRequest {
		var jsonErrBody apiErrorMessage
		if err := json.NewDecoder(resp.Body).Decode(&jsonErrBody); err != nil {
			return fmt.Errorf("failed to parse error response - %s", err)
		}

		return fmt.Errorf(jsonErrBody.Error.Message)
	}

	bb := &bytes.Buffer{}
	_, err = bb.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s submission) printResult() {
	msg := `

    Your solution has been submitted successfully.
    %s
`
	suffix := "View it at:\n\n    "
	if s.metadata.AutoApprove && s.metadata.Team == "" {
		suffix = "You can complete the exercise and unlock the next core exercise at:\n"
	}
	fmt.Fprintf(Err, msg, suffix)
	fmt.Fprintf(Out, "    %s\n\n", s.metadata.URL)
}

func (s submission) validate() error {
	if s.metadata.ID == "" {
		return errors.New("id is empty")
	}
	if len(s.documents) == 0 {
		return errors.New("documents is empty")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(submitCmd)
}

type apiErrorMessage struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}
