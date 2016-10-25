package api

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNewIteration_CaseSensitive(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(path, "..", "..", "fixtures", "iteration")

	testCases := []map[string][]string{
		{
			"file": []string{filepath.Join(dir, "python", "leap", "one.py")},
		},
		{
			"file": []string{filepath.Join(dir, "Python", "leap", "one.py")},
		},
		{
			"file": []string{filepath.Join(dir, "Python", "Leap", "one.py")},
		},
	}

	for _, testCase := range testCases {
		iter, err := NewIteration(dir, testCase["file"])
		if err != nil {
			t.Fatal(err)
		}

		if iter.TrackID != "python" {
			t.Errorf("Expected language to be python, was %s", iter.TrackID)
		}
		if iter.Problem != "leap" {
			t.Errorf("Expected problem to be leap, was %s", iter.Problem)
		}
	}
}

func TestNewIteration(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(path, "..", "..", "fixtures", "iteration")

	files := []string{
		filepath.Join(dir, "python", "leap", "one.py"),
		filepath.Join(dir, "python", "leap", "two.py"),
		filepath.Join(dir, "python", "leap", "lib", "three.py"),
		filepath.Join(dir, "python", "leap", "utf16le.py"),
		filepath.Join(dir, "python", "leap", "utf16be.py"),
		filepath.Join(dir, "python", "leap", "long-utf8.py"),
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

	if len(iter.Solution) != 6 {
		t.Fatalf("Expected solution to have 6 files, had %d", len(iter.Solution))
	}

	expected := map[string]struct {
		prefix string
		suffix string
	}{
		"one.py": {prefix: "# one"},
		"two.py": {prefix: "# two"},
		filepath.Join("lib", "three.py"): {prefix: "# three"},
		"utf16le.py":                     {prefix: "# utf16le"},
		"utf16be.py":                     {prefix: "# utf16be"},
		"long-utf8.py":                   {prefix: "# The first 1024", suffix: "👍\n"},
	}

	for filename, code := range expected {
		if !utf8.ValidString(iter.Solution[filename]) {
			t.Errorf("Iteration content is not valid UTF-8 data: %s", iter.Solution[filename])
		}

		if !strings.HasPrefix(iter.Solution[filename], code.prefix) {
			t.Errorf("Expected %s to start with `%s', had `%s'", filename, code.prefix, iter.Solution[filename])
		}
		if !strings.HasSuffix(iter.Solution[filename], code.suffix) {
			t.Errorf("Expected %s to end with `%s', had `%s'", filename, code.suffix, iter.Solution[filename])
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
