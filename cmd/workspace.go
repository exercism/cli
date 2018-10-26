package cmd

import (
	"fmt"

	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var checkTrack, checkExercise bool
var treeLevel int = 1
type caseInsensitive struct {
	values []string
}
// workspaceCmd outputs the path to the person's workspace directory.
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"w"},
	Short:   "Print out the path to your Exercism workspace.",
	Long: `Print out the path to your Exercism workspace.

This command can be used for scripting, or it can be combined with shell
commands to take you to your workspace.

For example you can run:

    cd $(exercism workspace)

On Windows, this will work only with Powershell, however you would
need to be on the same drive as your workspace directory. Otherwise
nothing will happen.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		

		var workspaceDir = v.GetString("workspace")
		if checkTrack {
			treeLevel = 1
			treePrint(workspaceDir)
		} else if checkExercise {
			treeLevel = 2
			treePrint(workspaceDir)
		} else {
			fmt.Fprintf(Out, "%s\n", workspaceDir)
		}
		return nil
	},
}

func treePrint(path string) (int) {
	path, err := filepath.Abs(path)
	if err != nil {
		return 0
	}
	_ = nodeVisit(path, "", 0)
	return 0
}

func nodeVisit(path, indent string, levelCheck int) (int) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	fmt.Fprintf(Out, "%s\n", fi.Name())
	if !fi.IsDir() {
		return 0
	}
	dir, err := os.Open(path)
	if err != nil {
		return 0
	}
	names, err := dir.Readdirnames(-1)
	_ = dir.Close() // safe to ignore this error.
	if err != nil {
		return 0
	}
	names = removeHidden(names)
	sort.Sort(caseInsensitive{names})
	if levelCheck == treeLevel{
		return 1
	}
	add := "│   "
	for i, name := range names {
		if i == len(names)-1 {
			fmt.Fprintf(Out, indent + "└── ")
			add = "    "
		} else {
			fmt.Fprintf(Out, indent + "├── ")
		}
		nodeVisit(filepath.Join(path, name), indent+add, levelCheck+1)
	}
	return 0
}

func removeHidden(files []string) []string {
	var clean []string
	for _, f := range files {
		if f[0] != '.' && !strings.Contains(f, ".go") && f != ".exercism" && !strings.Contains(f, ".json") {
			clean = append(clean, f)
		}
	}
	return clean
}
func (ci caseInsensitive) Len() int {
	return len(ci.values)
}

func (ci caseInsensitive) Less(i, j int) bool {
	return strings.ToLower(ci.values[i]) < strings.ToLower(ci.values[j])
}

func (ci caseInsensitive) Swap(i, j int) {
	ci.values[i], ci.values[j] = ci.values[j], ci.values[i]
}

func init() {
	RootCmd.AddCommand(workspaceCmd)
	workspaceCmd.Flags().BoolVarP(&checkTrack, "track", "t", false, "Shows current track")
	workspaceCmd.Flags().BoolVarP(&checkExercise, "exercise", "e", false, "Shows current exercise")
}
