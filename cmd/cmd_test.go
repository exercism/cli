package cmd

import "testing"

func TestIsTest(t *testing.T) {
	testCases := []struct {
		name   string
		isTest bool
	}{
		{
			name:   "problem/whatever-test.ext",
			isTest: true,
		},
		{
			name:   "problem/whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/whatever_test.spec.ext",
			isTest: true,
		},
		{
			name:   "problem/WhateverTest.ext",
			isTest: true,
		},
		{
			name:   "problem/Whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/whatever_test.ext",
			isTest: true,
		},
		{
			name:   "problem/whatever.ext",
			isTest: false,
		},
		{
			name:   "problem/test.ext",
			isTest: true,
		},
		{
			name:   "problem/Whatever.t", // perl
			isTest: true,
		},
		{
			name:   "whatever_spec.ext", // lua
			isTest: true,
		},
	}

	for _, tt := range testCases {
		if isTest(tt.name) != tt.isTest {
			t.Fatalf("Expected isTest(%s) to be %t", tt.name, tt.isTest)
		}
	}
}

func TestIsREADME(t *testing.T) {
	testCases := []struct {
		name     string
		isREADME bool
	}{
		{
			name:     "problem/README.md",
			isREADME: true,
		},
		{
			name:     "problem/README",
			isREADME: true,
		},
		{
			name:     "problem/README.txt",
			isREADME: true,
		},
		{
			name:     "problem/some_problem.py",
			isREADME: false,
		},
		{
			name:     "problem/readme.lua",
			isREADME: false,
		},
		{
			name:     "problem/readme_spec.lua",
			isREADME: false,
		},
	}

	for _, tt := range testCases {
		if isREADME(tt.name) != tt.isREADME {
			t.Fatalf("Expected isREADME(%s) to be %t", tt.name, tt.isREADME)
		}
	}
}
