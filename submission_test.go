package main

import (
	"testing"
)

var isTestTests = []struct {
	filename string
	expected bool
}{
	{"bob_test.rb", true},
	{"bob.spec.js", true},
	{"bob_test.exs", true},
	{"bob_test.clj", true},
	{"bob_test.py", true},
	{"bob_test.go", true},
	{"bob_test.hs", true},
	{"bob.rb", false},
}

func TestIsTest(t *testing.T) {
	for _, test := range isTestTests {
		result := IsTest(test.filename)
		if test.expected != result {
			t.Errorf("Filename [%s] should be a test file but is not.", test.filename)
		}
	}
}
