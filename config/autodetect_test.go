package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	language    string
	expectation []string
}{
	{"c", []string{"main.c", "main.h"}},
}

func TestAutoDetect(t *testing.T) {
	exerciseDir, err := ioutil.TempDir("", "test_autodetect")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(exerciseDir)

	srcDir := filepath.Join(exerciseDir, "src")
	os.Mkdir(srcDir, 0755)

	content := []byte("Very valid syntax.")
	// c
	files := []string{"bob.c", "bob_c.c", "BobC.c"}
	for i, fn := range files {
		fullPath := filepath.Join(srcDir, fn)
		ioutil.WriteFile(fullPath, content, os.ModePerm)
		files[i] = fullPath
	}

	solutionFiles, _ := FindSolutions(defaultTrackGlobs, "c", exerciseDir)
	assert.ElementsMatch(t, files, solutionFiles)
	// go
	// java
}
