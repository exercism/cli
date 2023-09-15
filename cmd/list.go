package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// This is a map of all the valid tracks on Exercism. It's here only because
// I don't want to have to make an API call to get the list of tracks.
var validTracks = map[string]bool{
	"8th":          true,
	"abap":         true,
	"awk":          true,
	"ballerina":    true,
	"bash":         true,
	"c":            true,
	"csharp":       true,
	"cpp":          true,
	"cfml":         true,
	"cobol":        true,
	"clojure":      true,
	"coffeescript": true,
	"common":       true,
	"crystal":      true,
	"d":            true,
	"dart":         true,
	"delphi":       true,
	"elixir":       true,
	"elm":          true,
	"emacs":        true,
	"erlang":       true,
	"fsharp":       true,
	"fortran":      true,
	"gleam":        true,
	"go":           true,
	"groovy":       true,
	"haskell":      true,
	"java":         true,
	"javascript":   true,
	"julia":        true,
	"kotlin":       true,
	"lfe":          true,
	"lua":          true,
	"mips":         true,
	"nim":          true,
	"ocaml":        true,
	"objective-c":  true,
	"php":          true,
	"plsql":        true,
	"perl":         true,
	"pharo":        true,
	"powershell":   true,
	"prolog":       true,
	"purescript":   true,
	"python":       true,
	"r":            true,
	"racket":       true,
	"raku":         true,
	"reasonml":     true,
	"red":          true,
	"ruby":         true,
	"rust":         true,
	"scala":        true,
	"scheme":       true,
	"standard":     true,
	"swift":        true,
	"tcl":          true,
	"typescript":   true,
	"unison":       true,
	"v":            true,
	"vim":          true,
	"visual":       true,
	"webassembly":  true,
	"wren":         true,
	"zig":          true,
	"jq":           true,
	"x86-64":       true,
}

func validateDifficulty(cmd *cobra.Command) error {
	level, err := cmd.Flags().GetString("level")
	if err != nil {
		return err
	}

	if level == "" {
		return nil
	}

	switch level {
	case "easy", "medium", "hard":
		return nil
	default:
		return errors.New("invalid level value. Please use: easy, medium or hard")
	}
}

func validateTrack(cmd *cobra.Command) error {
	track, err := cmd.Flags().GetString("track")
	if err != nil {
		return err
	}

	if track == "" {
		return nil
	}

	if !validTracks[track] {
		return errors.New("invalid track name. Please double check the track name")
	}
	return nil
}

func validateFlags(cmd *cobra.Command) error {
	if err := validateTrack(cmd); err != nil {
		return err
	}
	if err := validateDifficulty(cmd); err != nil {
		return err
	}
	return nil
}

type list struct {
	track      string
	level      string
	all        bool
	isDownload bool
	exercises  *allExercises
	tracks     *allTracks
}

func (l *list) set(flags *pflag.FlagSet) error {
	var err error
	l.track, err = flags.GetString("track")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return err
	}
	l.all, err = flags.GetBool("all")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return err
	}
	l.level, err = flags.GetString("level")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return err
	}
	l.isDownload, err = flags.GetBool("download")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return err
	}
	return nil
}

type track struct {
	Slug          string    `json:"slug,omitempty"`
	Title         string    `json:"title,omitempty"`
	Course        bool      `json:"course,omitempty"`
	NumConcepts   int       `json:"num_concepts,omitempty"`
	NumExercises  int       `json:"num_exercises,omitempty"`
	WebURL        string    `json:"web_url,omitempty"`
	IconURL       string    `json:"icon_url,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	LastTouchedAt time.Time `json:"last_touched_at,omitempty"`
	IsNew         bool      `json:"is_new,omitempty"`
	Links         struct {
		Self      string `json:"self,omitempty"`
		Exercises string `json:"exercises,omitempty"`
		Concepts  string `json:"concepts,omitempty"`
	} `json:"links,omitempty"`
	IsJoined              bool `json:"is_joined,omitempty"`
	NumLearntConcepts     int  `json:"num_learnt_concepts,omitempty"`
	NumCompletedExercises int  `json:"num_completed_exercises,omitempty"`
	NumSolutions          int  `json:"num_solutions,omitempty"`
	HasNotifications      bool `json:"has_notifications,omitempty"`
}

type allTracks struct {
	Tracks []track `json:"tracks,omitempty"`
}

type exercise struct {
	Slug          string `json:"slug,omitempty"`
	Type          string `json:"type,omitempty"`
	Title         string `json:"title,omitempty"`
	IconURL       string `json:"icon_url,omitempty"`
	Difficulty    string `json:"difficulty,omitempty"`
	Blurb         string `json:"blurb,omitempty"`
	IsExternal    bool   `json:"is_external,omitempty"`
	IsUnlocked    bool   `json:"is_unlocked,omitempty"`
	IsRecommended bool   `json:"is_recommended,omitempty"`
	Links         struct {
		Self string `json:"self,omitempty"`
	} `json:"links,omitempty"`
}

type allExercises struct {
	Exercises []exercise `json:"exercises,omitempty"`
}

// This simple map is used to sort the exercises by difficulty level.
// The custom sort order is: easy, medium, hard (ascending order).
var difficultyWeights = map[string]int{
	"easy":   1,
	"medium": 2,
	"hard":   3,
}

// ByDifficulty Custom sorter for slices of exercises within the allExercises struct
type ByDifficulty []exercise

func (bd ByDifficulty) Len() int      { return len(bd) }
func (bd ByDifficulty) Swap(i, j int) { bd[i], bd[j] = bd[j], bd[i] }
func (bd ByDifficulty) Less(i, j int) bool {
	return difficultyWeights[bd[i].Difficulty] < difficultyWeights[bd[j].Difficulty]
}

type exerciseFile struct {
	fileName    string
	exerciseDir string
	sourceURL   string
	slug        string
}

func (at allTracks) printer(showAll bool) {
	var counter int
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Track",
		"Joined",
		"Exercises",
		"Completed",
		"Progress",
		"Last Touched",
		"Track Link",
	})
	for _, track := range at.Tracks {
		joined := func() string {
			if track.IsJoined {
				return "Yes"
			}
			return "No"
		}
		if showAll || track.IsJoined {
			t.AppendRow(
				table.Row{
					track.Title,
					joined(),
					track.NumExercises,
					track.NumCompletedExercises,
					trackProgress(track.NumExercises, track.NumCompletedExercises),
					formatTime(track.LastTouchedAt),
					track.WebURL,
				})
			counter++
		}
	}
	t.AppendFooter(table.Row{"Total tracks", counter})
	t.Render()
}

func (ae allExercises) printer(difficultyLevel string) {
	var counter int
	var unlockedCounter int

	// Use custom sort order for the exercises before printing.
	sort.Sort(ByDifficulty(ae.Exercises))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Exercise",
		"Difficulty",
		"Exercise Link",
	})

	for _, exercise := range ae.Exercises {
		if exercise.IsUnlocked && (difficultyLevel == "" || exercise.Difficulty == difficultyLevel) {
			// This is a hack to check if the exercise is unlocked;
			// If the total number of unlocked exercises is 0, then we don't
			// want to print anything.
			// TODO: Find a better way to do this. Or maybe don't?
			unlockedCounter++
			t.AppendRow(
				table.Row{
					exercise.Title,
					exercise.Difficulty,
					fmt.Sprintf("https://exercism.org%s", exercise.Links.Self),
				})
			counter++
		}
	}
	if unlockedCounter == 0 {
		track := strings.TrimSuffix(
			strings.SplitAfter(ae.Exercises[0].Links.Self, "/")[2],
			"/",
		)
		fmt.Printf(
			"You don't have any %s exercises unlocked for %s track yet.\n",
			difficultyLevel,
			track,
		)
		return
	}
	t.AppendFooter(table.Row{"Total exercises:", counter})
	t.Render()
}

func (ae allExercises) prepareDownload(track string, difficultyLevel string) []string {
	var solutionLinks []string
	apiURL := "https://api.exercism.io/v1"
	for _, exercise := range ae.Exercises {
		if exercise.IsUnlocked && (difficultyLevel == "" || exercise.Difficulty == difficultyLevel) {
			solutionLinks = append(
				solutionLinks,
				fmt.Sprintf(
					"%s/solutions/latest?exercise_id=%s&track_id=%s",
					apiURL,
					exercise.Slug,
					track,
				),
			)
		}
	}
	return solutionLinks
}

func trackProgress(total int, completed int) string {
	if completed == 0 {
		return "0%"
	}
	percentage := float64(completed) / float64(total) * 100
	return fmt.Sprintf("%.0f%%", percentage)
}

func formatTime(at time.Time) string {
	if at.IsZero() {
		return "Never"
	}
	return at.Format("2006-01-02")
}

func runList(cfg config.Config, flags *pflag.FlagSet) error {
	l, err := newList(flags, cfg.UserViperConfig)
	if err != nil {
		return err
	}

	token := cfg.UserViperConfig.GetString("token")
	apibaseurl := cfg.UserViperConfig.GetString("apibaseurl")
	workspace := cfg.UserViperConfig.GetString("workspace")

	if l.track != "" && l.level != "" || l.isDownload {
		l.exercises.printer(l.level)
		if l.isDownload {
			solutionLinks := l.exercises.prepareDownload(l.track, l.level)
			totalLinks := len(solutionLinks)

			if totalLinks == 0 {
				// 0 means that there are no unlocked exercises
				// for the given track and difficulty level. Bail out.
				// Note: The message that's printed comes from the printer() method.
				return nil
			}

			if readyToDownload(totalLinks) {
				// Create the client here and pass it on to other download methods.
				client, err := api.NewClient(token, apibaseurl)
				if err != nil {
					return err
				}
				download(client, getExerciseFiles(collectExercises(client, solutionLinks), workspace))
			}
			return nil
		}
		return nil
	}

	if l.track != "" {
		l.exercises.printer("")
		return nil
	}

	if l.all {
		l.tracks.printer(true)
		return nil
	}

	l.tracks.printer(false)
	return nil
}

func download(client *api.Client, files []exerciseFile) {
	totalFiles := len(files)
	allTargetDirs := make(map[string]bool)
	bar := progressbar.Default(
		int64(totalFiles),
		fmt.Sprintf("Downloading %d files: ", totalFiles),
	)

	for _, file := range files {
		allTargetDirs[file.exerciseDir] = true

		resp, err := client.MakeRequest(file.sourceURL, true)
		if err != nil {
			// TODO: Export or handle the error better.
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		}(resp.Body)

		// Don't bother with empty files.
		if resp.Header.Get("Content-Length") == "0" {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "error: %s\n", resp.Status)
			os.Exit(1)
		}

		targetDir := filepath.Join(file.exerciseDir, filepath.Dir(file.fileName))
		createExerciseDir(targetDir)

		writeToFile(file.exerciseDir, file.fileName, resp.Body)
		_ = bar.Add(1)
	}
	printSuccess(totalFiles, allTargetDirs)
}

func writeToFile(dir string, fileName string, body io.ReadCloser) {
	targetFile, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	defer func(targetFile *os.File) {
		err := targetFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}(targetFile)

	_, err = io.Copy(targetFile, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func printSuccess(totalFiles int, allTargetDirs map[string]bool) {
	fmt.Printf(
		"\nSuccessfully downloaded %d files to the following directories:\n",
		totalFiles,
	)
	for dir := range allTargetDirs {
		fmt.Println(dir)
	}
	fmt.Println("Happy coding!")
}

func getExerciseFiles(solutions []ExerciseSolution, workspace string) []exerciseFile {
	var files []exerciseFile
	for _, solution := range solutions {
		for _, file := range solution.Solution.Files {
			ef := exerciseFile{
				fileName:    file,
				exerciseDir: getExerciseDir(solution, workspace),
				sourceURL:   fmt.Sprintf("%s%s", solution.Solution.FileDownloadBaseURL, file),
				slug:        solution.Solution.Exercise.ID,
			}
			files = append(files, ef)
		}
	}
	return files
}

func createExerciseDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}
}

func getExerciseDir(exercise ExerciseSolution, workspace string) string {
	metadata := exercise.GetSolutionMetadata()
	return metadata.Exercise(workspace).MetadataDir()
}

func collectExercises(client *api.Client, solutionLinks []string) []ExerciseSolution {
	solutions := make([]ExerciseSolution, 0, len(solutionLinks))
	fmt.Println("Collecting exercises...")
	for _, link := range solutionLinks {
		resp, err := client.MakeRequest(link, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "error: %s\n", resp.Status)
			os.Exit(1)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}

		solution := ExerciseSolution{}
		if err := json.Unmarshal(body, &solution); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		solutions = append(solutions, solution)
	}
	return solutions
}

func readyToDownload(count int) bool {
	var response string
	fmt.Printf("Do you want to download %d exercises? (y/n) ", count)
	fmt.Scanln(&response)
	if response == "y" || response == "Y" || response == "yes" || response == "Yes" {
		return true
	}
	return false
}

func newList(flags *pflag.FlagSet, usrCfg *viper.Viper) (*list, error) {
	var err error
	l := &list{}
	token := usrCfg.GetString("token")
	apibaseurl := usrCfg.GetString("apibaseurl")

	if err = l.set(flags); err != nil {
		return nil, err
	}

	// I'm using the experimental API to get the list of tracks and exercises.
	// The experimental API is not documented, and it's not guaranteed to be
	// stable. It's also not guaranteed to be available at all times.
	// The reason I'm using it is that the v1 API doesn't provide the
	// information that I need.
	// I reverse engineered the experimental API by looking at the source code
	// of the Exercism website.
	experimentalAPI := strings.Replace(apibaseurl, "v1", "v2", 1)
	url := fmt.Sprintf("%s/tracks", experimentalAPI)
	if l.track != "" {
		url = fmt.Sprintf("%s/tracks/%s/exercises", experimentalAPI, l.track)
	}

	client, err := api.NewClient(token, url)
	if err != nil {
		return nil, err
	}

	resp, err := client.MakeRequest(url, true)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if l.track != "" {
		l.exercises = &allExercises{}
		err = json.Unmarshal(body, l.exercises)
		if err != nil {
			return nil, err
		}
	} else {
		l.tracks = &allTracks{}
		err = json.Unmarshal(body, l.tracks)
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tracks",
	Long: `List tracks.

This command lists all the available tracks on Exercism.

It also lists all the exercises for a given track and difficulty level.

By default, it lists only the tracks that you have joined.

To list all the tracks, use the --all flag.

Finally, you can download all the exercises for a given track and
difficulty level using the --download flag.

Examples:

	exercism list // lists all the tracks that you have joined
	exercism list --all // lists all the tracks available on Exercism
	exercism list --track=go // lists all the exercises you're eligible for on the go track
	exercism list --track=go --level=easy // lists all the easy exercises on the go track
	exercism list --track=go --level=easy --download // downloads all the easy exercises on the go track
`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlags(cmd)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := LoadUserConfig()
		return runList(cfg, cmd.Flags())
	},
}

func setupListFlags(flags *pflag.FlagSet) {
	flags.StringP("track", "t", "", "track to list exercises for")
	flags.BoolP("all", "a", false, "list all tracks")
	flags.StringP("level", "l", "", "level to list exercises for")
	flags.BoolP("download", "d", false, "download the exercises")
}

func init() {
	RootCmd.AddCommand(listCmd)
	setupListFlags(listCmd.Flags())
}
