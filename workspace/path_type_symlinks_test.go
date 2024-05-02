//go:build !windows

package workspace

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectPathTypeSymlink(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "detect-path-type")

	testCases := []detectPathTestCase{
		{
			desc: "symlinked dir",
			path: filepath.Join(root, "symlinked-dir"),
			pt:   TypeDir,
		},
		{
			desc: "symlinked file",
			path: filepath.Join(root, "symlinked-file.txt"),
			pt:   TypeFile,
		},
	}

	for _, tc := range testCases {
		testDetectPathType(t, tc)
	}
}
