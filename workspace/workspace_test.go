package workspace

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateErrors(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	ws := New(filepath.Join(root, "workspace"))

	tests := []struct {
		desc, arg string
		errFn     func(error) bool
	}{
		{
			desc:  "absolute path outside of workspace",
			arg:   filepath.Join(root, "equipment", "bat"),
			errFn: IsNotInWorkspace,
		},
		{
			desc:  "absolute path in workspace not found",
			arg:   filepath.Join(ws.Dir, "creatures", "pig"),
			errFn: IsNotExist,
		},
		{
			desc:  "relative path is outside of workspace",
			arg:   filepath.Join("..", "fixtures", "locate-exercise", "equipment", "bat"),
			errFn: IsNotInWorkspace,
		},
		{
			desc:  "relative path in workspace not found",
			arg:   filepath.Join("..", "fixtures", "locate-exercise", "workspace", "creatures", "pig"),
			errFn: IsNotExist,
		},
		{
			desc:  "exercise name not found in workspace",
			arg:   "pig",
			errFn: IsNotExist,
		},
	}

	for _, test := range tests {
		_, err := ws.Locate(test.arg)
		assert.True(t, test.errFn(err), fmt.Sprintf("test: %s (arg: %s), %#v", test.desc, test.arg, err))
	}
}

func TestLocate(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	wsPrimary := New(filepath.Join(root, "workspace"))
	wsSymbolic := New(filepath.Join(root, "symlinked-workspace"))

	tests := []struct {
		desc      string
		workspace Workspace
		in        string
		out       []string
	}{
		{
			desc:      "find absolute path within workspace",
			workspace: wsPrimary,
			in:        filepath.Join(wsPrimary.Dir, "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find absolute path within symlinked workspace",
			workspace: wsSymbolic,
			in:        filepath.Join(wsSymbolic.Dir, "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find relative path within workspace",
			workspace: wsPrimary,
			in:        filepath.Join("..", "fixtures", "locate-exercise", "workspace", "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in default location",
			workspace: wsPrimary,
			in:        "horse",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in a symlinked workspace",
			workspace: wsSymbolic,
			in:        "horse",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in a subtree",
			workspace: wsPrimary,
			in:        "fly",
			out:       []string{filepath.Join(wsPrimary.Dir, "friends", "alice", "creatures", "fly")},
		},
		{
			desc:      "don't be confused by a file named the same as an exercise",
			workspace: wsPrimary,
			in:        "duck",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "duck")},
		},
		{
			desc:      "don't be confused by a symlinked file named the same as an exercise",
			workspace: wsPrimary,
			in:        "date",
			out:       []string{filepath.Join(wsPrimary.Dir, "actions", "date")},
		},
		{
			desc:      "find all the exercises with the same name",
			workspace: wsPrimary,
			in:        "bat",
			out: []string{
				filepath.Join(wsPrimary.Dir, "creatures", "bat"),
				filepath.Join(wsPrimary.Dir, "friends", "alice", "creatures", "bat"),
			},
		},
		{
			desc:      "find copies of exercise with suffix",
			workspace: wsPrimary,
			in:        "crane",
			out: []string{
				filepath.Join(wsPrimary.Dir, "creatures", "crane"),
				filepath.Join(wsPrimary.Dir, "creatures", "crane-2"),
			},
		},
		{
			desc:      "find exercises that are symlinks",
			workspace: wsPrimary,
			in:        "squash",
			out: []string{
				filepath.Join(wsPrimary.Dir, "..", "food", "squash"),
				filepath.Join(wsPrimary.Dir, "actions", "squash"),
			},
		},
	}

	for _, test := range tests {
		dirs, err := test.workspace.Locate(test.in)

		sort.Strings(dirs)
		sort.Strings(test.out)

		assert.NoError(t, err, test.desc)
		assert.Equal(t, test.out, dirs, test.desc)
	}
}
