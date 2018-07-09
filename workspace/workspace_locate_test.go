// +build !windows

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

	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := ws.Locate(tc.arg)
			assert.True(t, tc.errFn(err), fmt.Sprintf("test: %s (arg: %s), %#v", tc.desc, tc.arg, err))
		})
	}
}

type locateTestCase struct {
	desc      string
	workspace Workspace
	in        string
	out       []string
}

func TestLocate(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	wsPrimary := New(filepath.Join(root, "workspace"))

	testCases := []locateTestCase{
		{
			desc:      "find absolute path within workspace",
			workspace: wsPrimary,
			in:        filepath.Join(wsPrimary.Dir, "creatures", "horse"),
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
	}

	testLocate(testCases, t)
}

func testLocate(testCases []locateTestCase, t *testing.T) {
	for _, tc := range testCases {
		dirs, err := tc.workspace.Locate(tc.in)

		sort.Strings(dirs)
		sort.Strings(tc.out)

		assert.NoError(t, err, tc.desc)
		assert.Equal(t, tc.out, dirs, tc.desc)
	}
}
