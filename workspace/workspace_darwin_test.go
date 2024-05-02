package workspace

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExerciseDir_case_insensitive(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "solution-dir")
	// configuration file was set with "workspace" - the directory that exists
	configured := Workspace{Dir: filepath.Join(root, "workspace")}
	// user changes into directory with "bad" case - "Workspace"
	userPath := strings.Replace(configured.Dir, "workspace", "Workspace", 1)

	_, err := configured.ExerciseDir(filepath.Join(userPath, "exercise", "file.txt"))

	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("not in workspace: directory location may be case sensitive: "+
		"workspace directory: %s, submit path: %s/exercise/file.txt", configured.Dir, userPath), err.Error())
}
