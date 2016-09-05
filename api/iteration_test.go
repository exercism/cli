package api

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNewIteration(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(path, "..", "..", "fixtures", "iteration")

	files := []string{
		filepath.Join(dir, "python", "leap", "one.py"),
		filepath.Join(dir, "python", "leap", "two.py"),
		filepath.Join(dir, "python", "leap", "lib", "three.py"),
		filepath.Join(dir, "python", "leap", "utf16le.py"),
		filepath.Join(dir, "python", "leap", "utf16be.py"),
	}

	iter, err := NewIteration(dir, files)
	if err != nil {
		t.Fatal(err)
	}

	if iter.TrackID != "python" {
		t.Errorf("Expected language to be python, was %s", iter.TrackID)
	}
	if iter.Problem != "leap" {
		t.Errorf("Expected problem to be leap, was %s", iter.Problem)
	}

	if len(iter.Solution) != 5 {
		t.Fatalf("Expected solution to have 5 files, had %d", len(iter.Solution))
	}

	expected := map[string]string{
		"one.py": "# one",
		"two.py": "# two",
		filepath.Join("lib", "three.py"): "# three",
		"utf16le.py":                     "# utf16le",
		"utf16be.py":                     "# utf16be",
	}

	for filename, code := range expected {
		if !utf8.ValidString(iter.Solution[filename]) {
			t.Errorf("Iteration content is not valid UTF-8 data: %s", iter.Solution[filename])
		}

		if !strings.HasPrefix(iter.Solution[filename], code) {
			t.Errorf("Expected %s to contain `%s', had `%s'", filename, code, iter.Solution[filename])
		}
	}
}

func TestIterationValidFile(t *testing.T) {
	testCases := []struct {
		file string
		ok   bool
	}{
		{
			file: "/Users/me/exercism/ruby/bob/totally/fine/deep/path/src/bob.rb",
			ok:   true,
		},
		{
			file: "/Users/me/exercism/ruby/bob/bob.rb",
			ok:   true,
		},
		{
			file: "/users/me/exercism/ruby/bob/bob.rb",
			ok:   true,
		},
		{
			file: "/Users/me/bob.rb",
			ok:   false,
		},
		{
			file: "/tmp/bob.rb",
			ok:   false,
		},
	}

	for _, tt := range testCases {
		iter := &Iteration{
			Dir: "/Users/me/exercism",
		}
		ok := iter.isValidFilepath(tt.file)
		if ok && !tt.ok {
			t.Errorf("Expected %s to be invalid.", tt.file)
		}
	}
}
