package api

import (
	"path/filepath"
	"runtime"
	"testing"
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
		if iter.Exercise != "leap" {
			t.Errorf("Expected exercise to be leap, was %s", iter.Exercise)
		}
	}
}
