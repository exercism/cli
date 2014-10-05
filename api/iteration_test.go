package api

import "testing"

func TestIdentify(t *testing.T) {
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
			file: "/Users/me/exercism/bob.rb",
			ok:   false,
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
			File: tt.file,
			Dir:  "/Users/me/exercism",
		}
		err := iter.Identify()
		if !tt.ok && err == nil {
			t.Errorf("Expected %s to fail.", tt.file)
		}

		if tt.ok && !(iter.Language == "ruby" && iter.Problem == "bob") {
			t.Errorf("Language: %s, Problem: %s\nPath: %s\n", iter.Language, iter.Problem, tt.file)
		}
	}
}
