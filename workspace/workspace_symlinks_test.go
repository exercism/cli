// +build !windows

package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocateSymlinks(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "locate-exercise")

	wsSymbolic, err := New(filepath.Join(root, "symlinked-workspace"))
	assert.NoError(t, err)
	wsPrimary, err := New(filepath.Join(root, "workspace"))
	assert.NoError(t, err)

	testCases := []locateTestCase{
		{
			desc:      "find absolute path within symlinked workspace",
			workspace: wsSymbolic,
			in:        filepath.Join(wsSymbolic.Dir, "creatures", "horse"),
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "find by name in a symlinked workspace",
			workspace: wsSymbolic,
			in:        "horse",
			out:       []string{filepath.Join(wsPrimary.Dir, "creatures", "horse")},
		},
		{
			desc:      "don't be confused by a symlinked file named the same as an exercise",
			workspace: wsPrimary,
			in:        "date",
			out:       []string{filepath.Join(wsPrimary.Dir, "actions", "date")},
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

	testLocate(testCases, t)
}
