// +build windows

package workspace

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectPathType(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "detect-path-type")

	tests := []struct {
		desc string
		path string
		pt   PathType
	}{
		{
			desc: "absolute dir",
			path: filepath.Join(root, "a-dir"),
			pt:   TypeDir,
		},
		{
			desc: "relative dir",
			path: filepath.Join("..", "fixtures", "detect-path-type", "a-dir"),
			pt:   TypeDir,
		},
		{
			desc: "absolute file",
			path: filepath.Join(root, "a-file.txt"),
			pt:   TypeFile,
		},
		{
			desc: "relative file",
			path: filepath.Join("..", "fixtures", "detect-path-type", "a-file.txt"),
			pt:   TypeFile,
		},
		{
			desc: "symlinked file",
			path: filepath.Join(root, "symlinked-file.txt"),
			pt:   TypeFile,
		},
		{
			desc: "exercise ID",
			path: "a-file",
			pt:   TypeExerciseID,
		},
	}

	for _, test := range tests {
		pt, err := DetectPathType(test.path)
		assert.NoError(t, err, test.desc)
		assert.Equal(t, test.pt, pt, test.desc)
	}
}
func TestWindowsSymlinkedDir(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "detect-path-type")
	aDir := filepath.Join(root, "a-dir")
	symlink := filepath.Join(root, "symlinked-dir-windows")

	err := os.Symlink(aDir, "symlinked-dir-windows")
	// defer os.RemoveAll(symlink)

	pt, err := DetectPathType(symlink)
	assert.NoError(t, err, "windows symlinked dir")
	assert.Equal(t, TypeDir, pt, "windows symlinked dir")
}
