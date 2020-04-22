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

func setupTestDirectiory(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_autodetect")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)
	os.Chdir(tmpdir)

	srcDir := filepath.Join(tmpdir, "src")
	os.Mkdir(srcDir, 0755)

	content := []byte("Very valid syntax.")
	// c
	files := []string{"src/bob.c", "src/bob_c.c", "src/BobC.c"}
	for _, fn := range files {
		ioutil.WriteFile(fn, content, os.ModePerm)
	}
	solutionFiles, _ := findSolutions("c")
	assert.ElementsMatch(t, files, solutionFiles)
	// go
	// java
}

func TestAutoDetect(t *testing.T) {
	setupTestDirectiory(t)
}
