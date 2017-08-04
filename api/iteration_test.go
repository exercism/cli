package api

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFollowSymlink(t *testing.T) {
	// TODO: when we figure out where to deal with solutions
	// make sure we handle symlinks.
	t.Skip()
}

func TestHandleEncoding(t *testing.T) {
	// TODO: when we figure out where to deal with solutions
	// make sure we handle encoding properly.
	t.Skip()
	_, path, _, _ := runtime.Caller(0)
	dir := filepath.Join(path, "..", "..", "fixtures")

	_ = []string{
		filepath.Join(dir, "utf16le.py"),
		filepath.Join(dir, "utf16be.py"),
		filepath.Join(dir, "long-utf8.py"),
	}

	expected := map[string]struct {
		prefix string
		suffix string
	}{
		"utf16le.py":   {prefix: "# utf16le"},
		"utf16be.py":   {prefix: "# utf16be"},
		"long-utf8.py": {prefix: "# The first 1024", suffix: "👍\n"},
	}

	var solutions map[string]string
	for filename, code := range expected {
		if !utf8.ValidString(solutions[filename]) {
			t.Errorf("Iteration content is not valid UTF-8 data: %s", solutions[filename])
		}

		if !strings.HasPrefix(solutions[filename], code.prefix) {
			t.Errorf("Expected %s to start with `%s', had `%s'", filename, code.prefix, solutions[filename])
		}
		if !strings.HasSuffix(solutions[filename], code.suffix) {
			t.Errorf("Expected %s to end with `%s', had `%s'", filename, code.suffix, solutions[filename])
		}
	}
}
