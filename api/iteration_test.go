package api

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewIteration(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(path, "..", "..", "fixtures", "iteration")

	files := []string{
		filepath.Join(dir, "python", "leap", "one.py"),
		filepath.Join(dir, "python", "leap", "two.py"),
		filepath.Join(dir, "python", "leap", "lib", "three.py"),
	}

	iter, err := NewIteration(dir, files)
	if err != nil {
		t.Fatal(err)
	}

	if iter.Language != "python" {
		t.Errorf("Expected language to be python, was %s", iter.Language)
	}
	if iter.Problem != "leap" {
		t.Errorf("Expected problem to be leap, was %s", iter.Problem)
	}

	if len(iter.Solution) != 3 {
		t.Fatalf("Expected solution to have 3 files, had %d", len(iter.Solution))
	}

	expected := map[string]string{
		"one.py":       "# one\n",
		"two.py":       "# two\n",
		"lib/three.py": "# three\n",
	}
	for filename, code := range expected {
		if iter.Solution[filename] != code {
			t.Errorf("Expected %s to contain %s, had %s", filename, code, iter.Solution[filename])
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
