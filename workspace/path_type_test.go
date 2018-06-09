package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type detectPathTestCase struct {
	desc string
	path string
	pt   PathType
}

func TestDetectPathType(t *testing.T) {
	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "detect-path-type")

	testCases := []detectPathTestCase{
		detectPathTestCase{
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
			desc: "exercise ID",
			path: "a-file",
			pt:   TypeExerciseID,
		},
	}

	for _, tc := range testCases {
		testDetectPathType(t, tc)
	}
}

func testDetectPathType(t *testing.T, tc detectPathTestCase) {
	pt, err := DetectPathType(tc.path)
	assert.NoError(t, err, tc.desc)
	assert.Equal(t, tc.pt, pt, tc.desc)
}
