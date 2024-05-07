package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	Use:     "submit [<FILE> ...]",
	Aliases: []string{"s"},
	Short:   "Submit your solution to an exercise.",
	Long: `Submit your solution to an Exercism exercise.

    Call the command with the list of files you want to submit.
    If you omit the list of files, the CLI will submit the
    default solution files for the exercise.
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

		if len(args) == 0 {
			files, err := getExerciseSolutionFiles(".")
			if err != nil {
				return err
			}
			args = files
		}

		return runSubmit(cfg, cmd.Flags(), args)
	},
}

func runSubmit(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	if err := validateUserConfig(cfg.UserViperConfig); err != nil {
		return err
	}

	ctx := newSubmitCmdContext(cfg.UserViperConfig, flags)

	if err := ctx.validator.filesExistAndNotADir(args); err != nil {
		return err
	}

	submitPaths, err := ctx.evaluatedSymlinks(args)
	if err != nil {
		return err
	}

	submitPaths = ctx.removeDuplicatePaths(submitPaths)

	if err = ctx.validator.filesBelongToSameExercise(submitPaths); err != nil {
		return err
	}

	exercise, err := ctx.exercise(submitPaths[0])
	if err != nil {
		return err
	}

	if err = ctx.migrateLegacyMetadata(exercise); err != nil {
		return err
	}

	if err = ctx.validator.fileSizesWithinMax(submitPaths); err != nil {
		return err
	}

	documents, err := ctx.documents(submitPaths, exercise)
	if err != nil {
		return err
	}

	if err = ctx.validator.submissionNotEmpty(documents); err != nil {
		return err
	}

	metadata, err := ctx.metadata(exercise)
	if err != nil {
		return err
	}

	if err := ctx.validator.metadataMatchesExercise(metadata, exercise); err != nil {
		return err
	}

	if err := ctx.validator.isRequestor(metadata); err != nil {
		return err
	}

	if err := ctx.submit(metadata, documents); err != nil {
		return err
	}

	ctx.printResult(metadata)
	return nil
}

func getExerciseSolutionFiles(baseDir string) ([]string, error) {
	v := viper.New()
	v.AddConfigPath(filepath.Join(baseDir, ".exercism"))
	v.SetConfigName("config")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.New("no files to submit")
	}
	solutionFiles := v.GetStringSlice("files.solution")
	if len(solutionFiles) == 0 {
		return nil, errors.New("no files to submit")
	}

	return solutionFiles, nil
}

type submitCmdContext struct {
	usrCfg    *viper.Viper
	flags     *pflag.FlagSet
	validator submitValidator
}

func newSubmitCmdContext(usrCfg *viper.Viper, flags *pflag.FlagSet) *submitCmdContext {
	return &submitCmdContext{
		usrCfg:    usrCfg,
		flags:     flags,
		validator: submitValidator{usrCfg: usrCfg},
	}
}

// evaluatedSymlinks returns the submit paths with evaluated symlinks.
func (s *submitCmdContext) evaluatedSymlinks(submitPaths []string) ([]string, error) {
	evalSymlinkSubmitPaths := make([]string, 0, len(submitPaths))
	for _, path := range submitPaths {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		src, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}
		evalSymlinkSubmitPaths = append(evalSymlinkSubmitPaths, src)
	}
	return evalSymlinkSubmitPaths, nil
}

func (s *submitCmdContext) removeDuplicatePaths(submitPaths []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(submitPaths))

	for _, val := range submitPaths {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}

// exercise creates an exercise using one of the submitted filepaths.
// This assumes prior verification that submit paths belong to the same exercise.
func (s *submitCmdContext) exercise(aSubmitPath string) (workspace.Exercise, error) {
	ws, err := workspace.New(s.usrCfg.GetString("workspace"))
	if err != nil {
		return workspace.Exercise{}, err
	}

	dir, err := ws.ExerciseDir(aSubmitPath)
	if err != nil {
		return workspace.Exercise{}, err
	}
	return workspace.NewExerciseFromDir(dir), nil
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

// documents builds the documents that get submitted.
// Empty files are skipped, printing a warning.
func (s *submitCmdContext) documents(submitPaths []string, exercise workspace.Exercise) ([]workspace.Document, error) {
	docs := make([]workspace.Document, 0, len(submitPaths))
	for _, file := range submitPaths {
		// Don't submit empty files
		info, err := os.Stat(file)
		if err != nil {
			return nil, err
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
	return docs, nil
}

func (s *submitCmdContext) metadata(exercise workspace.Exercise) (*workspace.ExerciseMetadata, error) {
	metadata, err := workspace.NewExerciseMetadata(exercise.Filepath())
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

// submit submits the documents to the Exercism API.
func (s *submitCmdContext) submit(metadata *workspace.ExerciseMetadata, docs []workspace.Document) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, doc := range docs {
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

	client, err := api.NewClient(s.usrCfg.GetString("token"), s.usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/solutions/%s", s.usrCfg.GetString("apibaseurl"), metadata.ID)
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

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodedAPIError(resp)
	}

	bb := &bytes.Buffer{}
	_, err = bb.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s *submitCmdContext) printResult(metadata *workspace.ExerciseMetadata) {
	msg := `

    Your solution has been submitted successfully.
    %s
`
	suffix := "View it at:\n\n    "
	if metadata.AutoApprove && metadata.Team == "" {
		suffix = "You can complete the exercise and unlock the next core exercise at:\n"
	}
	fmt.Fprintf(Err, msg, suffix)
	fmt.Fprintf(Out, "    %s\n\n", metadata.URL)
}

// submitValidator contains the validation rules for a submission.
type submitValidator struct {
	usrCfg *viper.Viper
}

// filesExistAndNotADir checks that each file exists and is not a directory.
func (s submitValidator) filesExistAndNotADir(submitPaths []string) error {
	for _, path := range submitPaths {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				msg := `

    The file you are trying to submit cannot be found.

        %s

        `
				return fmt.Errorf(msg, path)
			}
			return err
		}
		if info.IsDir() {
			msg := `

    You are submitting a directory, which is not currently supported.

        %s

    Please change into the directory and provide the path to the file(s) you wish to submit

        %s submit FILENAME

            `
			return fmt.Errorf(msg, path, BinaryName)
		}
	}
	return nil
}

// filesBelongToSameExercise checks that each file belongs to the same exercise.
func (s submitValidator) filesBelongToSameExercise(submitPaths []string) error {
	ws, err := workspace.New(s.usrCfg.GetString("workspace"))
	if err != nil {
		return err
	}

	var exerciseDir string
	for _, f := range submitPaths {
		dir, err := ws.ExerciseDir(f)
		if err != nil {
			if workspace.IsMissingMetadata(err) {
				return errors.New(msgMissingMetadata)
			}
			return err
		}
		if exerciseDir != "" && dir != exerciseDir {
			msg := `

    You are submitting files belonging to different solutions.
    Please submit the files for one solution at a time.

        `
			return errors.New(msg)
		}
		exerciseDir = dir
	}
	return nil
}

// fileSizesWithinMax checks that each file does not exceed the max allowed size.
func (s submitValidator) fileSizesWithinMax(submitPaths []string) error {
	for _, file := range submitPaths {
		info, err := os.Stat(file)
		if err != nil {
			return err
		}
		const maxFileSize int64 = 65535
		if info.Size() >= maxFileSize {
			msg := `

      The submitted file '%s' is larger than the max allowed file size of %d bytes.
      Please reduce the size of the file and try again.

         `
			return fmt.Errorf(msg, file, maxFileSize)
		}
	}
	return nil
}

// submissionNotEmpty checks that there is at least one file to submit.
func (s submitValidator) submissionNotEmpty(docs []workspace.Document) error {
	if len(docs) == 0 {
		msg := `

    No files found to submit.

        `
		return errors.New(msg)
	}
	return nil
}

// metadataMatchesExercise checks that the metadata refers to the exercise being submitted.
func (s submitValidator) metadataMatchesExercise(metadata *workspace.ExerciseMetadata, exercise workspace.Exercise) error {
	if metadata.ExerciseSlug != exercise.Slug {
		// TODO: error msg should suggest running future doctor command
		msg := `

    The exercise directory does not match exercise slug in metadata:

        expected '%[1]s' but got '%[2]s'

    Please rename the directory '%[1]s' to '%[2]s' and try again.

        `
		return fmt.Errorf(msg, exercise.Slug, metadata.ExerciseSlug)
	}
	return nil
}

// isRequestor checks that the submission requestor is listed as the author in the metadata.
func (s submitValidator) isRequestor(metadata *workspace.ExerciseMetadata) error {
	if !metadata.IsRequester {
		msg := `

    The solution you are submitting is not connected to your account.
    Please re-download the exercise to make sure it has the data it needs.

        %s download --exercise=%s --track=%s

        `
		return fmt.Errorf(msg, BinaryName, metadata.ExerciseSlug, metadata.Track)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(submitCmd)
}
