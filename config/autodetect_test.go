package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutoDetect(t *testing.T) {
	exerciseDir, err := ioutil.TempDir("", "test_autodetect")
	assert.NoError(t, err)

	// Reset the test dir
	origDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(origDir)

	os.Chdir(exerciseDir)
	defer os.RemoveAll(exerciseDir)

	var testCases = []struct {
		Track    string
		Files    []string
		Expected []string
	}{
		{"c",
			[]string{"src/bob.c", "src/bob_c.h", "src/johnny.d"},
			[]string{
				filepath.Join(exerciseDir, "src/bob.c"),
				filepath.Join(exerciseDir, "src/bob_c.h"),
			},
		},
		{
			"bash",
			[]string{"bob.sh", "bob_test.sh"},
			[]string{
				filepath.Join(exerciseDir, "bob.sh"),
			},
		},
	}

	srcDir := filepath.Join(exerciseDir, "src")
	os.Mkdir(srcDir, 0755)

	// Arbitrary files
	err = createFiles([]string{"readme.md"}, exerciseDir)
	assert.NoError(t, err)

	for _, tc := range testCases {
		err = createFiles(tc.Files, exerciseDir)
		assert.NoError(t, err)

		solutionFiles, err := FindSolutions(defaultTrackGlobs, tc.Track, exerciseDir)
		t.Logf("Running test for track %s", tc.Track)
		assert.NoError(t, err)
		assert.ElementsMatch(t, tc.Expected, solutionFiles)
	}

}

func createFiles(files []string, dir string) error {
	content := []byte("Very valid syntax.")
	for _, fn := range files {
		fullPath := filepath.Join(dir, fn)
		err := ioutil.WriteFile(fullPath, content, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
